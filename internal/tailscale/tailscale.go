// Package tailscale is everything grat asks of Tailscale, and nothing more.
//
// The rest of grat never calls the foreign tool directly. It takes the Client
// interface, so the commands can be tested on a machine that has no Tailscale, and
// so a different way of publishing a local service would replace this package
// rather than every caller.
package tailscale

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Funnel is one local path published to the public internet.
type Funnel struct {
	// Path is the single path that goes public. It starts with a slash, and it is
	// the whole of what leaves the machine.
	Path string
	// PublicPort is the port the funnel listens on. Tailscale accepts 443, 8443 and
	// 10000, and grat rejects anything else when it reads the configuration.
	PublicPort int
	// Target is the local address the funnel forwards to, such as
	// http://localhost:4001.
	Target string
}

// SameAs reports whether two funnels are the same publication.
//
// All three fields decide it, and the target is the one that matters most. A
// path and a port say which slot is taken; only the target says which service is
// behind it. Comparing the first two alone made every service sharing the
// default path look published whilst one funnel existed, so grat status showed
// one address on every row.
func (funnel Funnel) SameAs(other Funnel) bool {
	return funnel.Path == other.Path &&
		funnel.PublicPort == other.PublicPort &&
		funnel.Target == other.Target
}

// IsAmong reports whether this funnel is one of those.
func (funnel Funnel) IsAmong(published []Funnel) bool {
	for _, candidate := range published {
		if funnel.SameAs(candidate) {
			return true
		}
	}
	return false
}

// PublicURL returns the address the outside world reaches this funnel at.
// hostname is the machine's name inside the tailnet, with or without the trailing
// dot that Tailscale reports. The default port is left out of the URL, because a
// payment provider stores this address and a shorter one is easier to check.
func (funnel Funnel) PublicURL(hostname string) string {
	host := strings.TrimSuffix(hostname, ".")
	if funnel.PublicPort != defaultPublicPort {
		host = host + ":" + strconv.Itoa(funnel.PublicPort)
	}
	return "https://" + host + funnel.Path
}

// defaultPublicPort is the port a funnel uses when the configuration names none.
// It matches the default in the config package, and both are the port Tailscale
// itself defaults to.
const defaultPublicPort = 443

// Client is the set of Tailscale operations grat performs.
//
// Implementations must leave alone anything grat did not open. Close therefore
// takes the same funnel that Open received, rather than clearing whatever is
// currently published.
type Client interface {
	// Open publishes funnel.Path at funnel.PublicPort, forwarding to funnel.Target.
	//
	// needsEnabling is called with an address when the tailnet has not enabled
	// Funnel, which is a permission only its owner can grant. The call happens
	// whilst Open is still waiting, because that address is what ends the wait.
	Open(ctx context.Context, funnel Funnel, needsEnabling func(address string)) error
	// Close withdraws exactly the given funnel and leaves every other one standing.
	Close(ctx context.Context, funnel Funnel) error
	// Funnels reports what is published right now.
	Funnels(ctx context.Context) ([]Funnel, error)
	// Hostname returns the machine's name inside the tailnet.
	Hostname(ctx context.Context) (string, error)
}

// ErrNotInstalled reports that no Tailscale command could be found. It is not a
// failure in itself: the caller decides whether to install Tailscale or to say
// what is missing.
type ErrNotInstalled struct {
	// Searched lists the places that were looked at, so the message can say where.
	Searched []string
}

func (err ErrNotInstalled) Error() string {
	if len(err.Searched) == 0 {
		return "tailscale is not installed"
	}
	return "tailscale is not installed; looked in " + strings.Join(err.Searched, ", ")
}

// ErrCommandFailed reports that the Tailscale command ran and refused. It carries
// what Tailscale said, so a caller can show a reason without the caller having to
// know anything about the tool.
type ErrCommandFailed struct {
	// Arguments is the command line grat used, without the executable.
	Arguments []string
	// Output is what Tailscale wrote, trimmed.
	Output string
	// Err is the underlying failure from running the process.
	Err error
}

func (err ErrCommandFailed) Error() string {
	if err.Output == "" {
		return fmt.Sprintf("tailscale %s failed: %v", strings.Join(err.Arguments, " "), err.Err)
	}
	return fmt.Sprintf("tailscale %s failed: %s", strings.Join(err.Arguments, " "), err.Output)
}

func (err ErrCommandFailed) Unwrap() error {
	return err.Err
}
