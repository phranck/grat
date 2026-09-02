package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phranck/grat/internal/config"
)

// TestAHealthProbeDoesNotFollowARedirectOffTheService is the fault this guards.
// The readiness client was the default one, which follows up to ten redirects
// wherever they point, so a development server answering its health path with a
// redirect had grat fetch whatever it named, from inside the developer's network.
func TestAHealthProbeDoesNotFollowARedirectOffTheService(t *testing.T) {
	t.Parallel()

	reached := make(chan string, 1)
	elsewhere := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		reached <- request.URL.Path
		writer.WriteHeader(http.StatusOK)
	}))
	defer elsewhere.Close()

	service := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, elsewhere.URL+"/secrets", http.StatusFound)
	}))
	defer service.Close()

	client := probeClient(t)
	response, err := client.Get(service.URL + "/health")
	if err == nil {
		_ = response.Body.Close()
		t.Fatal("the probe followed the redirect off the service")
	}
	if !strings.Contains(err.Error(), "redirected to") {
		t.Fatalf("error = %v, want the refusal to leave the service", err)
	}
	select {
	case path := <-reached:
		t.Fatalf("the probe reached the other host at %s", path)
	default:
	}
}

// TestAHealthProbeFollowsARedirectOnTheSameService keeps the bound from turning
// an ordinary redirect into a failure. A service moving its health path to
// another path of its own is still that service answering.
func TestAHealthProbeFollowsARedirectOnTheSameService(t *testing.T) {
	t.Parallel()

	service := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/health" {
			http.Redirect(writer, request, "/up", http.StatusFound)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
	defer service.Close()

	response, err := probeClient(t).Get(service.URL + "/health")
	if err != nil {
		t.Fatalf("Get() error = %v, want a redirect within the service to be followed", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
}

// probeClient is the client a readiness probe uses, taken from a normalized
// manager so the test cannot disagree with what runs.
func probeClient(t *testing.T) *http.Client {
	t.Helper()
	manager, err := Manager{
		Root: t.TempDir(),
		Config: config.Config{
			Version: 1,
			Project: config.Project{Name: "fixture"},
			Services: []config.Service{{
				Name: "frontend", Command: "npm run dev", Role: config.RoleFrontend,
				Port: 3000, Host: "localhost", HealthPath: "/health",
			}},
			Runtime: config.DefaultRuntime(),
		},
	}.normalized()
	if err != nil {
		t.Fatalf("normalized() error = %v", err)
	}
	return manager.httpClient()
}
