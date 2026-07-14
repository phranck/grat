package maintenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestUpdateDelegatesToHomebrewOnlyForOwnedExecutable(t *testing.T) {
	t.Parallel()

	prefix := filepath.Join(t.TempDir(), "Cellar", "grat", "1.0.0")
	executable := filepath.Join(prefix, "bin", "grat")
	if err := os.MkdirAll(filepath.Dir(executable), 0o700); err != nil {
		t.Fatalf("create formula executable directory: %v", err)
	}
	if err := os.WriteFile(executable, []byte("old"), 0o755); err != nil {
		t.Fatalf("write formula executable: %v", err)
	}
	commands := &fakeCommands{responses: map[string]commandResponse{
		commandKey("brew", "list", "--versions", HomebrewFormula): {output: []byte("grat 1.0.0\n")},
		commandKey("brew", "--prefix", HomebrewFormula):           {output: []byte(prefix + "\n")},
		commandKey("brew", "upgrade", HomebrewFormula):            {},
	}}
	service := Service{
		Executable:     func() (string, error) { return executable, nil },
		EvalSymlinks:   filepath.EvalSymlinks,
		Command:        commands.Run,
		CurrentVersion: func() string { return "v1.0.0" },
	}

	result, err := service.Update(context.Background())
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if !strings.Contains(result.Message, "Homebrew") {
		t.Fatalf("Update() message = %q, want Homebrew result", result.Message)
	}
	if !commands.called(commandKey("brew", "upgrade", HomebrewFormula)) {
		t.Fatalf("commands = %#v, want brew upgrade", commands.calls)
	}
}

func TestUpdateDoesNotDelegateToHomebrewForAnotherExecutable(t *testing.T) {
	t.Parallel()

	prefix := filepath.Join(t.TempDir(), "Cellar", "grat", "1.0.0")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("go installation"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	commands := &fakeCommands{responses: map[string]commandResponse{
		commandKey("brew", "list", "--versions", HomebrewFormula): {output: []byte("grat 1.0.0\n")},
		commandKey("brew", "--prefix", HomebrewFormula):           {output: []byte(prefix + "\n")},
	}}
	service := Service{
		Executable:     func() (string, error) { return executable, nil },
		EvalSymlinks:   filepath.EvalSymlinks,
		Command:        commands.Run,
		CurrentVersion: func() string { return "v1.0.0" },
		BuildInfo:      func() (string, string, bool) { return ModulePath, "v1.0.0", true },
	}

	result, err := service.Update(context.Background())
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if got, want := result.Message, "Run: go install github.com/phranck/grat/cmd/grat@latest"; got != want {
		t.Fatalf("Update() message = %q, want %q", got, want)
	}
	if commands.called(commandKey("brew", "upgrade", HomebrewFormula)) {
		t.Fatalf("commands = %#v, must not upgrade a non-Homebrew executable", commands.calls)
	}
}

func TestUpdateReplacesVerifiedReleaseForEverySupportedTarget(t *testing.T) {
	for _, target := range []struct{ goos, goarch string }{
		{goos: "darwin", goarch: "amd64"},
		{goos: "darwin", goarch: "arm64"},
		{goos: "linux", goarch: "amd64"},
		{goos: "linux", goarch: "arm64"},
	} {
		target := target
		t.Run(target.goos+"/"+target.goarch, func(t *testing.T) {
			t.Parallel()

			oldBinary := []byte("old release " + target.goos + "/" + target.goarch)
			newBinary := []byte("new release " + target.goos + "/" + target.goarch)
			executable := filepath.Join(t.TempDir(), "grat")
			if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
				t.Fatalf("write executable: %v", err)
			}
			server := newReleaseServer(t, target.goos, target.goarch, oldBinary, newBinary, false)
			defer server.Close()
			service := releaseService(executable, server.URL, target.goos, target.goarch)

			result, err := service.Update(context.Background())
			if err != nil {
				t.Fatalf("Update() error = %v", err)
			}
			if !strings.Contains(result.Message, "v1.0.1") {
				t.Fatalf("Update() message = %q, want latest version", result.Message)
			}
			got, err := os.ReadFile(executable)
			if err != nil {
				t.Fatalf("read replacement executable: %v", err)
			}
			if string(got) != string(newBinary) {
				t.Fatalf("replacement executable = %q, want %q", got, newBinary)
			}
		})
	}
}

func TestUpdateKeepsCurrentReleaseWhenChecksumDoesNotMatch(t *testing.T) {
	t.Parallel()

	oldBinary := []byte("old release")
	newBinary := []byte("new release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newReleaseServer(t, runtime.GOOS, runtime.GOARCH, oldBinary, newBinary, true)
	defer server.Close()

	_, err := releaseService(executable, server.URL, runtime.GOOS, runtime.GOARCH).Update(context.Background())
	if err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("Update() error = %v, want checksum failure", err)
	}
	got, readErr := os.ReadFile(executable)
	if readErr != nil {
		t.Fatalf("read executable after failed update: %v", readErr)
	}
	if string(got) != string(oldBinary) {
		t.Fatalf("executable changed after failed update: got %q, want %q", got, oldBinary)
	}
}

func TestUpdateDoesNotReplaceBinaryWhenAtomicRenameFails(t *testing.T) {
	t.Parallel()

	oldBinary := []byte("old release")
	newBinary := []byte("new release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newReleaseServer(t, runtime.GOOS, runtime.GOARCH, oldBinary, newBinary, false)
	defer server.Close()
	service := releaseService(executable, server.URL, runtime.GOOS, runtime.GOARCH)
	service.Rename = func(string, string) error { return errors.New("rename denied") }

	if _, err := service.Update(context.Background()); err == nil {
		t.Fatal("Update() error = nil, want rename failure")
	}
	got, readErr := os.ReadFile(executable)
	if readErr != nil {
		t.Fatalf("read executable after failed update: %v", readErr)
	}
	if string(got) != string(oldBinary) {
		t.Fatalf("executable changed after failed update: got %q, want %q", got, oldBinary)
	}
}

func TestReleaseChecksumRejectsOversizedDocument(t *testing.T) {
	for _, streamed := range []bool{false, true} {
		streamed := streamed
		t.Run(fmt.Sprintf("streamed=%t", streamed), func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				if streamed {
					writer.(http.Flusher).Flush()
				} else {
					writer.Header().Set("Content-Length", "33")
				}
				_, _ = writer.Write([]byte(strings.Repeat("x", 33)))
			}))
			defer server.Close()
			service := Service{
				ReleaseAPI:              server.URL,
				HTTPClient:              server.Client(),
				MaxReleaseDocumentBytes: 32,
			}
			value := release{Assets: []asset{{Name: "checksums.txt", BrowserDownloadURL: server.URL}}}

			_, err := service.releaseChecksum(context.Background(), value, "grat_v1.0.0_darwin_arm64")
			if err == nil || !strings.Contains(err.Error(), "exceeds") {
				t.Fatalf("releaseChecksum() error = %v, want size-limit failure", err)
			}
		})
	}
}

func TestReplaceVerifiedBinaryRejectsOversizedDownload(t *testing.T) {
	for _, streamed := range []bool{false, true} {
		streamed := streamed
		t.Run(fmt.Sprintf("streamed=%t", streamed), func(t *testing.T) {
			t.Parallel()

			payload := []byte("oversized release")
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				if streamed {
					writer.(http.Flusher).Flush()
				} else {
					writer.Header().Set("Content-Length", fmt.Sprint(len(payload)))
				}
				_, _ = writer.Write(payload)
			}))
			defer server.Close()
			executable := filepath.Join(t.TempDir(), "grat")
			original := []byte("current release")
			if err := os.WriteFile(executable, original, 0o755); err != nil {
				t.Fatalf("write executable: %v", err)
			}
			service := Service{
				ReleaseAPI:           server.URL,
				HTTPClient:           server.Client(),
				MaxReleaseAssetBytes: int64(len(payload) - 1),
			}

			err := service.replaceVerifiedBinary(context.Background(), executable, server.URL, digest(payload))
			if err == nil || !strings.Contains(err.Error(), "exceeds") {
				t.Fatalf("replaceVerifiedBinary() error = %v, want size-limit failure", err)
			}
			got, readErr := os.ReadFile(executable)
			if readErr != nil {
				t.Fatalf("read executable: %v", readErr)
			}
			if string(got) != string(original) {
				t.Fatalf("executable changed after oversized download: got %q, want %q", got, original)
			}
		})
	}
}

func releaseService(executable string, apiURL string, goos string, goarch string) Service {
	return Service{
		Executable:   func() (string, error) { return executable, nil },
		EvalSymlinks: filepath.EvalSymlinks,
		Command: func(context.Context, string, ...string) ([]byte, error) {
			return nil, errors.New("not installed by Homebrew")
		},
		CurrentVersion: func() string { return "v1.0.0" },
		ReleaseAPI:     apiURL,
		GOOS:           goos,
		GOARCH:         goarch,
		Rename:         os.Rename,
	}
}

func newReleaseServer(t *testing.T, goos string, goarch string, oldBinary []byte, newBinary []byte, corruptLatestChecksum bool) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/repos/phranck/grat/releases/tags/v1.0.0":
			writeReleaseJSON(t, writer, "v1.0.0", goos, goarch, "/current/checksums.txt", "/current/grat", oldBinary)
		case "/repos/phranck/grat/releases/latest":
			payload := newBinary
			if corruptLatestChecksum {
				payload = []byte("different bytes")
			}
			writeReleaseJSON(t, writer, "v1.0.1", goos, goarch, "/latest/checksums.txt", "/latest/grat", payload)
		case "/current/checksums.txt":
			_, _ = fmt.Fprintf(writer, "%s  grat_v1.0.0_%s_%s\n", digest(oldBinary), goos, goarch)
		case "/latest/checksums.txt":
			payload := newBinary
			if corruptLatestChecksum {
				payload = []byte("different bytes")
			}
			_, _ = fmt.Fprintf(writer, "%s  grat_v1.0.1_%s_%s\n", digest(payload), goos, goarch)
		case "/current/grat":
			_, _ = writer.Write(oldBinary)
		case "/latest/grat":
			_, _ = writer.Write(newBinary)
		default:
			http.NotFound(writer, request)
		}
	}))
	return server
}

func writeReleaseJSON(t *testing.T, writer http.ResponseWriter, tag string, goos string, goarch string, checksumPath string, binaryPath string, checksumBytes []byte) {
	t.Helper()
	base := "http://" + writer.Header().Get("Host")
	if base == "http://" {
		base = ""
	}
	_, err := fmt.Fprintf(writer, `{"tag_name":%q,"assets":[{"name":"checksums.txt","browser_download_url":%q},{"name":%q,"browser_download_url":%q}]}`,
		tag,
		base+checksumPath,
		fmt.Sprintf("grat_%s_%s_%s", tag, goos, goarch),
		base+binaryPath,
	)
	if err != nil {
		t.Fatalf("write release response: %v", err)
	}
	_ = checksumBytes
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

type commandResponse struct {
	output []byte
	err    error
}

type fakeCommands struct {
	responses map[string]commandResponse
	calls     []string
}

func (commands *fakeCommands) Run(_ context.Context, name string, arguments ...string) ([]byte, error) {
	key := commandKey(name, arguments...)
	commands.calls = append(commands.calls, key)
	response, exists := commands.responses[key]
	if !exists {
		return nil, errors.New("unexpected command: " + key)
	}
	return response.output, response.err
}

func (commands *fakeCommands) called(wanted string) bool {
	for _, call := range commands.calls {
		if call == wanted {
			return true
		}
	}
	return false
}

func commandKey(name string, arguments ...string) string {
	return name + "\x00" + strings.Join(arguments, "\x00")
}
