package maintenance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestUpdateRefusesAnOlderLatestRelease is the fault this guards. grat took
// whatever GitHub marked as the latest release and compared it with the
// installed version for equality alone, so a release marked by mistake, or by
// somebody who reached the release listing, was installed however old it was.
// The checksum and the attestation prove where a binary came from and say
// nothing about its age.
func TestUpdateRefusesAnOlderLatestRelease(t *testing.T) {
	t.Parallel()

	installed := []byte("the installed release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, installed, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newOrderedReleaseServer(t, "v1.9.1", "v1.9.0", installed)
	defer server.Close()

	service := releaseService(executable, server, runtime.GOOS, runtime.GOARCH)
	service.CurrentVersion = func() string { return "v1.9.1" }

	_, err := service.Update(context.Background())
	if err == nil {
		t.Fatal("Update() error = nil, want a refusal to go backwards")
	}
	// Both versions, because the reader has to be able to see which way round
	// this went without looking anything up.
	for _, wanted := range []string{"v1.9.0", "v1.9.1"} {
		if !strings.Contains(err.Error(), wanted) {
			t.Fatalf("the refusal does not name %q: %v", wanted, err)
		}
	}
	got, readErr := os.ReadFile(executable)
	if readErr != nil {
		t.Fatalf("read executable after the refusal: %v", readErr)
	}
	if string(got) != string(installed) {
		t.Fatalf("the executable was replaced: got %q, want %q", got, installed)
	}
}

// TestUpdateReportsNothingToDoForTheSameRelease keeps the ordering from turning
// an ordinary up-to-date run into a refusal.
func TestUpdateReportsNothingToDoForTheSameRelease(t *testing.T) {
	t.Parallel()

	installed := []byte("the installed release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, installed, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newOrderedReleaseServer(t, "v1.9.1", "v1.9.1", installed)
	defer server.Close()

	service := releaseService(executable, server, runtime.GOOS, runtime.GOARCH)
	service.CurrentVersion = func() string { return "v1.9.1" }

	result, err := service.Update(context.Background())
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if !strings.Contains(result.Message, "v1.9.1") {
		t.Fatalf("Update() message = %q, want the installed version named", result.Message)
	}
}

// newOrderedReleaseServer answers with one release for the installed tag and one
// for latest, both carrying the same binary, so the only thing under test is
// which of the two versions is the newer.
func newOrderedReleaseServer(t *testing.T, installedTag string, latestTag string, binary []byte) *httptest.Server {
	t.Helper()
	goos, goarch := runtime.GOOS, runtime.GOARCH
	return httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/repos/phranck/grat/releases/tags/" + installedTag:
			writeReleaseJSON(t, writer, installedTag, goos, goarch, "/installed/checksums.txt", "/installed/grat", binary)
		case "/repos/phranck/grat/releases/latest":
			writeReleaseJSON(t, writer, latestTag, goos, goarch, "/latest/checksums.txt", "/latest/grat", binary)
		case "/installed/checksums.txt":
			_, _ = fmt.Fprintf(writer, "%s  grat_%s_%s_%s\n", digest(binary), installedTag, goos, goarch)
		case "/latest/checksums.txt":
			_, _ = fmt.Fprintf(writer, "%s  grat_%s_%s_%s\n", digest(binary), latestTag, goos, goarch)
		case "/installed/grat", "/latest/grat":
			_, _ = writer.Write(binary)
		default:
			http.NotFound(writer, request)
		}
	}))
}
