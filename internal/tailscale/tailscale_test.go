package tailscale

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestOpeningAFunnelUsesTheDocumentedArguments(t *testing.T) {
	t.Parallel()

	got := funnelArguments(Funnel{Path: "/api/webhooks/creem", PublicPort: 443, Target: "http://localhost:4001"}, false)
	want := []string{"funnel", "--yes", "--bg", "--https=443", "--set-path=/api/webhooks/creem", "http://localhost:4001"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("funnelArguments() = %q, want %q", got, want)
	}
}

func TestClosingAFunnelRepeatsEveryFlagAndDropsTheTarget(t *testing.T) {
	t.Parallel()

	got := funnelArguments(Funnel{Path: "/hook", PublicPort: 8443, Target: "http://localhost:4001"}, true)
	want := []string{"funnel", "--yes", "--https=8443", "--set-path=/hook", "off"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("funnelArguments() = %q, want %q", got, want)
	}
}

func TestThePublicURLLeavesOutTheDefaultPort(t *testing.T) {
	t.Parallel()

	funnel := Funnel{Path: "/api/webhooks/creem", PublicPort: 443}
	if got, want := funnel.PublicURL("fixture.tail1234.ts.net."), "https://fixture.tail1234.ts.net/api/webhooks/creem"; got != want {
		t.Fatalf("PublicURL() = %q, want %q", got, want)
	}
}

func TestThePublicURLNamesAnyOtherPort(t *testing.T) {
	t.Parallel()

	funnel := Funnel{Path: "/hook", PublicPort: 8443}
	if got, want := funnel.PublicURL("fixture.tail1234.ts.net"), "https://fixture.tail1234.ts.net:8443/hook"; got != want {
		t.Fatalf("PublicURL() = %q, want %q", got, want)
	}
}

func TestParsingTheServeConfigReportsOnlyPublicPaths(t *testing.T) {
	t.Parallel()

	output := []byte(`{
	  "AllowFunnel": {"fixture.tail1234.ts.net:443": true, "fixture.tail1234.ts.net:8443": false},
	  "Web": {
	    "fixture.tail1234.ts.net:443": {"Handlers": {"/api/webhooks/creem": {"Proxy": "http://127.0.0.1:4001"}}},
	    "fixture.tail1234.ts.net:8443": {"Handlers": {"/private": {"Proxy": "http://127.0.0.1:4002"}}}
	  }
	}`)

	got, err := parseFunnels(output)
	if err != nil {
		t.Fatalf("parseFunnels() error = %v", err)
	}
	want := []Funnel{{Path: "/api/webhooks/creem", PublicPort: 443, Target: "http://127.0.0.1:4001"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseFunnels() = %+v, want %+v", got, want)
	}
}

func TestParsingTheServeConfigReportsEveryPathOfOnePublicPort(t *testing.T) {
	t.Parallel()

	output := []byte(`{
	  "AllowFunnel": {"fixture.tail1234.ts.net:443": true},
	  "Web": {
	    "fixture.tail1234.ts.net:443": {"Handlers": {
	      "/first": {"Proxy": "http://127.0.0.1:4001"},
	      "/second": {"Proxy": "http://127.0.0.1:4002"}
	    }}
	  }
	}`)

	got, err := parseFunnels(output)
	if err != nil {
		t.Fatalf("parseFunnels() error = %v", err)
	}
	sort.Slice(got, func(left, right int) bool { return got[left].Path < got[right].Path })
	want := []Funnel{
		{Path: "/first", PublicPort: 443, Target: "http://127.0.0.1:4001"},
		{Path: "/second", PublicPort: 443, Target: "http://127.0.0.1:4002"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseFunnels() = %+v, want %+v", got, want)
	}
}

func TestParsingAnEmptyServeConfigReportsNothing(t *testing.T) {
	t.Parallel()

	for _, output := range []string{"", "null", "{}"} {
		got, err := parseFunnels([]byte(output))
		if err != nil {
			t.Fatalf("parseFunnels(%q) error = %v", output, err)
		}
		if len(got) != 0 {
			t.Fatalf("parseFunnels(%q) = %+v, want nothing published", output, got)
		}
	}
}

func TestParsingRejectsAHostPortWithoutAUsablePort(t *testing.T) {
	t.Parallel()

	output := []byte(`{"AllowFunnel": {"fixture.tail1234.ts.net": true}, "Web": {"fixture.tail1234.ts.net": {"Handlers": {}}}}`)
	if _, err := parseFunnels(output); err == nil {
		t.Fatal("parseFunnels() error = nil, want a refusal naming the missing port")
	}
}

func TestReadingTheStatusTakesTheMachineNameAndState(t *testing.T) {
	t.Parallel()

	value, err := parseStatus([]byte(`{"BackendState":"Running","Self":{"DNSName":"fixture.tail1234.ts.net."}}`))
	if err != nil {
		t.Fatalf("parseStatus() error = %v", err)
	}
	if value.BackendState != "Running" {
		t.Fatalf("BackendState = %q, want Running", value.BackendState)
	}
	if value.Self == nil || value.Self.DNSName != "fixture.tail1234.ts.net." {
		t.Fatalf("Self = %+v, want the reported DNS name", value.Self)
	}
}

func TestAMissingInstallationNamesWhereItLooked(t *testing.T) {
	t.Parallel()

	err := ErrNotInstalled{Searched: []string{"PATH", bundledExecutable}}
	if got := err.Error(); !strings.Contains(got, bundledExecutable) {
		t.Fatalf("Error() = %q, want the searched locations named", got)
	}
}
