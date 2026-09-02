// Package publish holds the rules that decide what of a project reaches the
// public internet, and which published funnel belongs to which service.
//
// It sits above internal/tailscale on purpose. The tailscale package is meant to
// be replaceable by another way of publishing a local service, so it must not
// learn grat's rules; and the rules themselves are answers a web interface needs
// as much as a terminal does, so nothing here takes a renderer, a flag set or an
// input stream. The command line parses what was typed, calls these functions and
// renders the result; a daemon calls the same functions from a request handler.
package publish

import (
	"errors"
	"fmt"
	"strings"

	"github.com/phranck/grat/internal/config"
	"github.com/phranck/grat/internal/tailscale"
	"github.com/phranck/grat/internal/textsafe"
)

// AllServices is the word that stands where a service name stands and selects
// every one of them. A word rather than a flag, because it reads as what it
// selects and sits where the thing it replaces sits.
const AllServices = "all"

// maxPathBytes bounds a path given on the command line, matching what the
// configuration allows for the same value.
const maxPathBytes = 2 << 10

// Publication is one service together with the funnel that publishes it.
type Publication struct {
	// Service is the configured service behind the funnel.
	Service config.Service
	// Funnel is the path, the public port and the local target that Tailscale is
	// asked to publish.
	Funnel tailscale.Funnel
}

// Whole reports whether this publication puts everything the service offers on
// the internet, which is what the root path means.
func (publication Publication) Whole() bool {
	return publication.Funnel.Path == config.DefaultExposePath
}

// Selection is what the publishing rule makes of the names a command was given.
type Selection struct {
	// Publications are the funnels to open, in the order the names were given.
	Publications []Publication
	// PassedOver names the services that all skipped because they name no path.
	// A command reports them, so nobody is left believing their project is
	// public when part of it is not.
	PassedOver []string
}

// Select resolves service names and an optional path into the funnels to open.
//
// A service is published only where a path says so, either as pathOverride from
// the command line or as path in its [services.expose] table. Publishing is the
// one thing grat does that reaches the internet, and a development server is not
// built for that: a debug toolbar shows itself to any request that appears to
// come from the machine itself, and funnel traffic does, because Tailscale dials
// the local target from there. So the whole of a service goes public only when
// somebody writes the root path down.
//
// The word all takes every service that carries an expose table and reports the
// rest through PassedOver. Naming a service explicitly that has no path is an
// error, because the name says what somebody meant.
func Select(value config.Config, names []string, pathOverride string) (Selection, error) {
	if pathOverride != "" {
		if err := ValidatePath(pathOverride); err != nil {
			return Selection{}, err
		}
		if len(names) > 1 || (len(names) == 1 && names[0] == AllServices) {
			return Selection{}, errors.New("--path names one path, so it takes one service")
		}
	}

	selection, err := selectPublications(value, names, pathOverride)
	if err != nil {
		return Selection{}, err
	}
	if err := refuseCollidingFunnels(selection.Publications); err != nil {
		return Selection{}, err
	}
	return selection, nil
}

// selectPublications applies the path rule to the named services, or to every
// service that carries an expose table where the name was all.
func selectPublications(value config.Config, names []string, pathOverride string) (Selection, error) {
	if len(names) == 1 && names[0] == AllServices {
		return selectEverythingWithAPath(value)
	}

	selection := Selection{Publications: make([]Publication, 0, len(names))}
	for _, name := range names {
		if name == AllServices {
			return Selection{}, errors.New("all selects every service, so it takes no other name beside it")
		}
		service, err := exposableService(value, name)
		if err != nil {
			return Selection{}, err
		}
		funnel, err := FunnelFor(service, pathOverride)
		if err != nil {
			return Selection{}, err
		}
		selection.Publications = append(selection.Publications, Publication{Service: service, Funnel: funnel})
	}
	return selection, nil
}

// selectEverythingWithAPath answers the word all: every service that says where
// it wants to be published, and the names of those that do not.
func selectEverythingWithAPath(value config.Config) (Selection, error) {
	selection := Selection{}
	addressable := 0
	for _, service := range value.Services {
		if service.Port == 0 {
			continue
		}
		addressable++
		funnel, err := FunnelFor(service, "")
		if err != nil {
			selection.PassedOver = append(selection.PassedOver, service.Name)
			continue
		}
		selection.Publications = append(selection.Publications, Publication{Service: service, Funnel: funnel})
	}
	if addressable == 0 {
		return Selection{}, errors.New("this project has no service with an address to publish")
	}
	if len(selection.Publications) == 0 {
		return Selection{}, errors.New(
			"no service of this project names a path to publish; give one a [services.expose] table with a path, or publish it once with grat expose NAME --path /some/path",
		)
	}
	return selection, nil
}

// Named resolves service names into the services behind them, without asking
// where they would be published.
//
// It answers the commands that close or report rather than open, since those
// find what is published by its target instead of deriving it from a path. The
// word all means every service that has an address at all.
func Named(value config.Config, names []string) ([]config.Service, error) {
	if len(names) == 1 && names[0] == AllServices {
		services := make([]config.Service, 0, len(value.Services))
		for _, service := range value.Services {
			if service.Port != 0 {
				services = append(services, service)
			}
		}
		if len(services) == 0 {
			return nil, errors.New("this project has no service with an address to publish")
		}
		return services, nil
	}

	services := make([]config.Service, 0, len(names))
	for _, name := range names {
		if name == AllServices {
			return nil, errors.New("all selects every service, so it takes no other name beside it")
		}
		service, err := exposableService(value, name)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

// FunnelFor turns a service into the funnel that publishes it, or says that
// nothing names a path for it.
//
// Two places can name that path, and the one given for this run wins over the
// one in the configuration:
//
//   - pathOverride, which is --path on the command line, for a single run
//   - a [services.expose] table, for a path that always applies
//
// Either of them may be the root path, and that is how a whole service is
// published: by writing it down rather than by leaving it out.
func FunnelFor(service config.Service, pathOverride string) (tailscale.Funnel, error) {
	path, publicPort := service.Exposure()
	if pathOverride != "" {
		path = pathOverride
	}
	if path == "" {
		return tailscale.Funnel{}, fmt.Errorf(
			"%s names no path to publish; add --path /some/path for this run, or a [services.expose] table with a path, and path = %q publishes all of it",
			service.Name, config.DefaultExposePath,
		)
	}
	return tailscale.Funnel{
		Path:       path,
		PublicPort: publicPort,
		Target:     TargetFor(service),
	}, nil
}

// TargetFor returns the local address a funnel forwards to for this service.
//
// It is what identifies a funnel as belonging to a service, because the path can
// be anything somebody typed whilst the target is the service's own address.
func TargetFor(service config.Service) string {
	return strings.TrimSuffix(service.URL(), "/")
}

// FunnelsFor returns the funnels among published that point at this service.
//
// Looking a funnel up by its target rather than by the path the configuration
// would derive is what makes a funnel opened with --path visible to the commands
// that report and close: the path was a choice made at the command line and is
// nowhere in the configuration, whilst the target is the service's address and
// is the same either way.
func FunnelsFor(service config.Service, published []tailscale.Funnel) []tailscale.Funnel {
	target := TargetFor(service)
	if target == "" {
		return nil
	}
	found := make([]tailscale.Funnel, 0, 1)
	for _, funnel := range published {
		if funnel.Target == target {
			found = append(found, funnel)
		}
	}
	return found
}

// ValidatePath checks a path given outside the configuration against the same
// rules the configuration applies to one written down in it.
func ValidatePath(path string) error {
	if len(path) > maxPathBytes {
		return fmt.Errorf("a path may be at most %d bytes long", maxPathBytes)
	}
	if textsafe.ContainsUnsafe(path) {
		return errors.New("a path must not contain control or Unicode format characters")
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("a path must begin with a slash, got %q", path)
	}
	return nil
}

// exposableService returns the named service. Every HTTP service can be
// published; a process-only worker cannot, because it has no address at all.
func exposableService(value config.Config, name string) (config.Service, error) {
	for _, service := range value.Services {
		if service.Name != name {
			continue
		}
		if service.Port == 0 {
			return config.Service{}, fmt.Errorf("%s is a process-only service and has no address to publish", name)
		}
		return service, nil
	}
	return config.Service{}, fmt.Errorf("unknown service %q", name)
}

// refuseCollidingFunnels reports two services that would take the same funnel.
//
// A funnel is its public port and its path, not the service behind it, so two
// services that name the same path both take the same slot and the second
// replaces the first. Publishing them one after another leaves the project with
// one public address and grat having said it opened two.
//
// It refuses before anything is published rather than after, because a partial
// publication leaves a project half public with no single command to undo it.
func refuseCollidingFunnels(publications []Publication) error {
	type slot struct {
		path string
		port int
	}
	taken := map[slot]string{}
	for _, publication := range publications {
		key := slot{path: publication.Funnel.Path, port: publication.Funnel.PublicPort}
		if first, exists := taken[key]; exists {
			return fmt.Errorf(
				"%s and %s would both be published at %s on port %d, and a funnel is that path and that port rather than the service behind it; give one of them its own path in a [services.expose] table",
				first, publication.Service.Name, key.path, key.port,
			)
		}
		taken[key] = publication.Service.Name
	}
	return nil
}
