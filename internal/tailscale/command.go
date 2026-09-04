package tailscale

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// bundledExecutable is where the macOS application keeps its command. That variant
// does not put it on the PATH, so grat looks here before concluding that Tailscale
// is missing.
const bundledExecutable = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"

// bundledCLIEnvironment makes the bundled binary behave as a command rather than
// opening the graphical application.
const bundledCLIEnvironment = "TAILSCALE_BE_CLI=1"

// maxCommandOutputBytes bounds what grat reads from the foreign tool, the way the
// config reader bounds a config file.
const maxCommandOutputBytes = 1 << 20

// CommandClient talks to Tailscale by running its command.
type CommandClient struct {
	// executable is the resolved path to the command.
	executable string
	// bundled reports whether executable is the one inside the macOS application,
	// which needs an environment variable to act as a command.
	bundled bool
}

// Locate finds the Tailscale command and returns a client bound to it.
//
// It looks on the PATH first and then inside the macOS application bundle. A
// missing command yields ErrNotInstalled naming both places, because that error is
// the starting point for installing Tailscale rather than a failure on its own.
func Locate() (CommandClient, error) {
	if path, err := exec.LookPath("tailscale"); err == nil {
		return CommandClient{executable: path}, nil
	}
	if info, err := os.Stat(bundledExecutable); err == nil && !info.IsDir() {
		return CommandClient{executable: bundledExecutable, bundled: true}, nil
	}
	return CommandClient{}, ErrNotInstalled{Searched: []string{"PATH", bundledExecutable}}
}

// NewCommandClient returns a client that runs the command at path. The setup step
// uses it after an installation, when it knows where the command landed.
func NewCommandClient(path string) CommandClient {
	return CommandClient{executable: path, bundled: path == bundledExecutable}
}

// Executable returns the resolved path, which the setup step reports.
func (client CommandClient) Executable() string {
	return client.executable
}

// funnelEnableMarker is what Tailscale prints when the tailnet has not enabled
// Funnel. It follows the line with an address and then waits, so the address has
// to be read out of the running command rather than after it.
const funnelEnableMarker = "https://login.tailscale.com/f/funnel"

// Open publishes one path.
//
// The funnel runs in the background, because it has to outlive the grat command
// that opened it: without --bg the tool holds the terminal and withdraws the
// funnel when it is interrupted. Prompts are answered in advance with --yes,
// since grat asks the one question that matters before it gets here.
//
// One thing --yes cannot answer is a tailnet that has not enabled Funnel at all.
// Tailscale then prints an address and waits for somebody to visit it, so the
// command neither fails nor returns. needsEnabling is called with that address
// when it appears, which is what lets the caller say so and open the page rather
// than leaving a reader in front of a command that looks stuck.
func (client CommandClient) Open(ctx context.Context, funnel Funnel, needsEnabling func(address string)) error {
	announced := false
	return client.stream(ctx, func(line string) {
		if announced || !strings.Contains(line, funnelEnableMarker) {
			return
		}
		announced = true
		if needsEnabling != nil {
			needsEnabling(strings.TrimSpace(line))
		}
	}, funnelArguments(funnel, false)...)
}

// Close withdraws one path. Tailscale requires the same flags the funnel was
// opened with, followed by off, so the arguments differ from Open only in that
// word and in the omitted target.
func (client CommandClient) Close(ctx context.Context, funnel Funnel) error {
	_, err := client.run(ctx, funnelArguments(funnel, true)...)
	return err
}

// funnelArguments builds the command line for opening or closing one funnel.
func funnelArguments(funnel Funnel, closing bool) []string {
	arguments := []string{"funnel", "--yes"}
	if !closing {
		arguments = append(arguments, "--bg")
	}
	arguments = append(arguments,
		"--https="+strconv.Itoa(funnel.PublicPort),
		"--set-path="+funnel.Path,
	)
	if closing {
		return append(arguments, "off")
	}
	return append(arguments, funnel.Target)
}

// Funnels reports what Tailscale currently publishes.
//
// The shape read here is the serve configuration: AllowFunnel says which host and
// port pairs are public, and Web lists the paths served below each of them. Only
// those two parts are read, so a field elsewhere in that structure does not
// concern grat.
func (client CommandClient) Funnels(ctx context.Context) ([]Funnel, error) {
	output, err := client.run(ctx, "funnel", "status", "--json")
	if err != nil {
		return nil, err
	}
	return parseFunnels(output)
}

// serveConfig is the part of Tailscale's serve configuration that grat reads.
type serveConfig struct {
	AllowFunnel map[string]bool            `json:"AllowFunnel"`
	Web         map[string]webServerConfig `json:"Web"`
}

type webServerConfig struct {
	Handlers map[string]webHandler `json:"Handlers"`
}

type webHandler struct {
	Proxy string `json:"Proxy"`
}

// parseFunnels turns the serve configuration into the funnels grat understands. A
// host and port pair that is not public is skipped, because a served path that is
// not funnelled never leaves the tailnet.
func parseFunnels(output []byte) ([]Funnel, error) {
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" || trimmed == "null" {
		return nil, nil
	}
	var config serveConfig
	if err := json.Unmarshal([]byte(trimmed), &config); err != nil {
		return nil, fmt.Errorf("read the published funnels: %w", err)
	}

	funnels := make([]Funnel, 0, len(config.Web))
	for hostPort, web := range config.Web {
		if !config.AllowFunnel[hostPort] {
			continue
		}
		port, err := portOf(hostPort)
		if err != nil {
			return nil, err
		}
		for path, handler := range web.Handlers {
			funnels = append(funnels, Funnel{Path: path, PublicPort: port, Target: handler.Proxy})
		}
	}
	return funnels, nil
}

// portOf extracts the port from a host and port pair such as
// machine.tail1234.ts.net:443.
func portOf(hostPort string) (int, error) {
	_, portText, found := strings.Cut(hostPort, ":")
	if !found {
		return 0, fmt.Errorf("read the published funnels: %q names no port", hostPort)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		return 0, fmt.Errorf("read the published funnels: %q names no usable port", hostPort)
	}
	return port, nil
}

// status is the part of Tailscale's status output that grat reads.
//
// AuthURL carries the sign-in address while a machine is signed out, so grat can
// open that page itself instead of reading it out of the tool's printed output.
type status struct {
	BackendState string `json:"BackendState"`
	AuthURL      string `json:"AuthURL"`
	Self         *struct {
		DNSName string `json:"DNSName"`
	} `json:"Self"`
}

// Hostname returns this machine's name inside the tailnet, without the trailing
// dot that Tailscale appends to the fully qualified name.
func (client CommandClient) Hostname(ctx context.Context) (string, error) {
	output, err := client.run(ctx, "status", "--json")
	if err != nil {
		return "", err
	}
	value, err := parseStatus(output)
	if err != nil {
		return "", err
	}
	if value.Self == nil || strings.TrimSpace(value.Self.DNSName) == "" {
		return "", errors.New("this machine has no name in the tailnet yet")
	}
	return strings.TrimSuffix(value.Self.DNSName, "."), nil
}

// parseStatus reads the status output into the fields grat uses.
func parseStatus(output []byte) (status, error) {
	var value status
	if err := json.Unmarshal(output, &value); err != nil {
		return status{}, fmt.Errorf("read the Tailscale status: %w", err)
	}
	return value, nil
}

// run executes one Tailscale command and returns its standard output. A refusal
// carries what Tailscale said, so the caller never has to read the tool's own
// error format.
func (client CommandClient) run(ctx context.Context, arguments ...string) ([]byte, error) {
	if client.executable == "" {
		return nil, ErrNotInstalled{}
	}
	// #nosec G204 -- the executable is a resolved absolute path and every argument
	// is built in this package from validated configuration.
	command := exec.CommandContext(ctx, client.executable, arguments...)
	if client.bundled {
		command.Env = append(os.Environ(), bundledCLIEnvironment)
	}

	var output, failure bytes.Buffer
	command.Stdout = &boundedWriter{writer: &output, remaining: maxCommandOutputBytes}
	command.Stderr = &boundedWriter{writer: &failure, remaining: maxCommandOutputBytes}
	if err := command.Run(); err != nil {
		return nil, ErrCommandFailed{Arguments: arguments, Output: strings.TrimSpace(failure.String()), Err: err}
	}
	return output.Bytes(), nil
}

// stream executes one Tailscale command and reports each line it prints whilst it
// is still running.
//
// run buffers everything and hands it over at the end, which is right for a
// command that answers a question. It is wrong for one that reports something
// and then waits, because the report is the reason for the wait and a reader
// sees nothing until the wait is over.
func (client CommandClient) stream(ctx context.Context, onLine func(string), arguments ...string) error {
	if client.executable == "" {
		return ErrNotInstalled{}
	}
	// #nosec G204 -- the executable is a resolved absolute path and every argument
	// is built in this package from validated configuration.
	command := exec.CommandContext(ctx, client.executable, arguments...)
	if client.bundled {
		command.Env = append(os.Environ(), bundledCLIEnvironment)
	}

	reader, writer := io.Pipe()
	command.Stdout = writer
	command.Stderr = writer

	var seen strings.Builder
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if seen.Len() < maxCommandOutputBytes {
				seen.WriteString(line)
				seen.WriteString("\n")
			}
			onLine(line)
		}
	}()

	err := command.Run()
	_ = writer.Close()
	<-done
	if err != nil {
		return ErrCommandFailed{Arguments: arguments, Output: strings.TrimSpace(seen.String()), Err: err}
	}
	return nil
}

// boundedWriter passes through at most remaining bytes and discards the rest, so a
// foreign process cannot make grat hold an unbounded amount of its output.
type boundedWriter struct {
	writer    io.Writer
	remaining int
}

func (bounded *boundedWriter) Write(data []byte) (int, error) {
	if bounded.remaining <= 0 {
		return len(data), nil
	}
	accepted := data
	if len(accepted) > bounded.remaining {
		accepted = accepted[:bounded.remaining]
	}
	written, err := bounded.writer.Write(accepted)
	bounded.remaining -= written
	if err != nil {
		return written, err
	}
	return len(data), nil
}
