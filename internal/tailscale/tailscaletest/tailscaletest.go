// Package tailscaletest provides a Tailscale client that records what it was asked
// to do, so commands can be tested on a machine without Tailscale.
package tailscaletest

import (
	"context"

	"github.com/phranck/grat/internal/tailscale"
)

// Client records every call and answers from fields the test sets. The zero value
// is a working client that publishes nothing and never fails.
type Client struct {
	// Opened and Closed record the funnels the caller asked for, in order.
	Opened []tailscale.Funnel
	Closed []tailscale.Funnel
	// Published is what Funnels reports.
	Published []tailscale.Funnel
	// Name is what Hostname reports. It defaults to a name shaped like a real one.
	Name string

	// The Err fields make the matching call fail, so a test can drive the failure
	// path without constructing a broken machine.
	OpenErr     error
	CloseErr    error
	FunnelsErr  error
	HostnameErr error
}

// Open records the funnel unless OpenErr is set.
func (client *Client) Open(_ context.Context, funnel tailscale.Funnel) error {
	if client.OpenErr != nil {
		return client.OpenErr
	}
	client.Opened = append(client.Opened, funnel)
	client.Published = append(client.Published, funnel)
	return nil
}

// Close records the funnel and removes it from what Funnels reports.
func (client *Client) Close(_ context.Context, funnel tailscale.Funnel) error {
	if client.CloseErr != nil {
		return client.CloseErr
	}
	client.Closed = append(client.Closed, funnel)
	remaining := make([]tailscale.Funnel, 0, len(client.Published))
	for _, published := range client.Published {
		if published.Path == funnel.Path && published.PublicPort == funnel.PublicPort {
			continue
		}
		remaining = append(remaining, published)
	}
	client.Published = remaining
	return nil
}

// Funnels reports what this client currently publishes.
func (client *Client) Funnels(context.Context) ([]tailscale.Funnel, error) {
	if client.FunnelsErr != nil {
		return nil, client.FunnelsErr
	}
	return client.Published, nil
}

// Hostname reports Name, or a default that looks like a real tailnet name.
func (client *Client) Hostname(context.Context) (string, error) {
	if client.HostnameErr != nil {
		return "", client.HostnameErr
	}
	if client.Name == "" {
		return "fixture.tail1234.ts.net", nil
	}
	return client.Name, nil
}
