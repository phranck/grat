package publish

import (
	"context"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/tailscale"
)

// Withdrawal is one public address that has been taken back, or that could not
// be.
type Withdrawal struct {
	// Service is the service the funnel forwarded to.
	Service config.Service
	// Funnel is what was published.
	Funnel tailscale.Funnel
	// Reopen is the command that puts the same address back. A funnel's address
	// is the tailnet name and the path, so reopening one gives the address it
	// had, and a webhook registered with a provider loses nothing by this.
	Reopen string
	// Err is why the funnel could not be closed, and is nil where it was.
	Err error
}

// WithdrawalObserver is told what became of each public address.
//
// It is a seam rather than a return value because the front ends differ in what
// they do with it: a terminal prints a line as each one closes, and a web
// interface updates a row. The rule itself neither prints nor asks.
type WithdrawalObserver interface {
	ObserveWithdrawal(withdrawal Withdrawal)
}

// Withdraw closes every published funnel that forwards to one of these services.
//
// A funnel outlives the service behind it, because it is configuration in
// Tailscale rather than a process. Left standing after a stop or a port change,
// it goes on forwarding to a local port, and whatever binds that port next is
// what answers the internet: another project's service after the numbers moved,
// or anything else that happens to take it.
//
// Nothing is asked here. The observer is told what closed and how to put it
// back, which is what the caller renders.
func Withdraw(ctx context.Context, client tailscale.Client, services []config.Service, observer WithdrawalObserver) error {
	published, err := client.Funnels(ctx)
	if err != nil {
		return err
	}
	if len(published) == 0 {
		return nil
	}
	for _, service := range services {
		for _, funnel := range FunnelsFor(service, published) {
			withdrawal := Withdrawal{
				Service: service,
				Funnel:  funnel,
				Reopen:  ReopenCommand(service, funnel),
				Err:     client.Close(ctx, funnel),
			}
			if observer != nil {
				observer.ObserveWithdrawal(withdrawal)
			}
		}
	}
	return nil
}

// ReopenCommand is the line that publishes this address again.
func ReopenCommand(service config.Service, funnel tailscale.Funnel) string {
	return "grat expose " + service.Name + " --path " + funnel.Path
}

// Open reports the funnels that are published for these services right now.
//
// It answers the commands that start something rather than stop it, so a
// service coming up under an address that was already open can say so instead
// of leaving somebody to find out from Tailscale.
func Open(ctx context.Context, client tailscale.Client, services []config.Service) ([]Publication, error) {
	published, err := client.Funnels(ctx)
	if err != nil {
		return nil, err
	}
	found := make([]Publication, 0, len(services))
	for _, service := range services {
		for _, funnel := range FunnelsFor(service, published) {
			found = append(found, Publication{Service: service, Funnel: funnel})
		}
	}
	return found, nil
}
