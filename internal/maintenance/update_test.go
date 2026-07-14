package maintenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
			service := releaseService(executable, server, target.goos, target.goarch)

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

	_, err := releaseService(executable, server, runtime.GOOS, runtime.GOARCH).Update(context.Background())
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

func TestUpdateKeepsCurrentReleaseWhenAttestationFails(t *testing.T) {
	t.Parallel()

	oldBinary := []byte("old release")
	newBinary := []byte("new release")
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, oldBinary, 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	server := newReleaseServer(t, runtime.GOOS, runtime.GOARCH, oldBinary, newBinary, false)
	defer server.Close()
	service := releaseService(executable, server, runtime.GOOS, runtime.GOARCH)
	service.VerifyAttestation = func(_ context.Context, path string, tag string) error {
		if tag == "v1.0.0" {
			return nil
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("attestation candidate is unavailable: %v", err)
		}
		return errors.New("attestation fixture rejected")
	}

	_, err := service.Update(context.Background())
	if err == nil || !strings.Contains(err.Error(), "attestation") {
		t.Fatalf("Update() error = %v, want attestation failure", err)
	}
	got, readErr := os.ReadFile(executable)
	if readErr != nil {
		t.Fatalf("read executable after failed update: %v", readErr)
	}
	if string(got) != string(oldBinary) {
		t.Fatalf("executable changed after failed attestation: got %q, want %q", got, oldBinary)
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
	service := releaseService(executable, server, runtime.GOOS, runtime.GOARCH)
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

			server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
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
			server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
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

			err := service.replaceVerifiedBinary(context.Background(), executable, server.URL, digest(payload), "v1.0.1")
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

func TestAssetURLRejectsUntrustedOrigins(t *testing.T) {
	t.Parallel()

	service := Service{ReleaseAPI: defaultReleaseAPI}
	for _, endpoint := range []string{
		"http://github.com/phranck/grat/releases/download/v1.0.0/grat_v1.0.0_darwin_arm64",
		"https://example.com/phranck/grat/releases/download/v1.0.0/grat_v1.0.0_darwin_arm64",
		"https://github.com/another/project/releases/download/v1.0.0/grat",
	} {
		if _, err := service.assetURL(endpoint); err == nil {
			t.Fatalf("assetURL(%q) accepted an untrusted release origin", endpoint)
		}
	}
}

func TestReleaseURLsAcceptOnlyCanonicalGitHubPaths(t *testing.T) {
	t.Parallel()

	service := Service{ReleaseAPI: defaultReleaseAPI}
	assetEndpoint := "https://github.com/phranck/grat/releases/download/v1.2.3/grat_v1.2.3_darwin_arm64"
	if got, err := service.assetURL(assetEndpoint); err != nil || got != assetEndpoint {
		t.Fatalf("assetURL() = (%q, %v), want canonical release asset", got, err)
	}
	redirect, err := url.Parse("https://release-assets.githubusercontent.com/github-production-release-asset/fixture")
	if err != nil {
		t.Fatalf("parse redirect fixture: %v", err)
	}
	if err := service.validateReleaseAssetURL(redirect, true); err != nil {
		t.Fatalf("validateReleaseAssetURL() rejected GitHub release storage: %v", err)
	}
	if _, err := (Service{ReleaseAPI: "http://api.github.com"}).endpoint(releaseAPIPathPrefix + "latest"); err == nil {
		t.Fatal("endpoint() accepted an insecure release API")
	}
}

func TestReleaseDownloadRejectsCrossOriginRedirect(t *testing.T) {
	t.Parallel()

	foreign := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte("foreign release"))
	}))
	defer foreign.Close()
	trusted := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, foreign.URL+"/asset", http.StatusFound)
	}))
	defer trusted.Close()
	executable := filepath.Join(t.TempDir(), "grat")
	if err := os.WriteFile(executable, []byte("current release"), 0o755); err != nil {
		t.Fatalf("write executable: %v", err)
	}
	service := Service{ReleaseAPI: trusted.URL, HTTPClient: trusted.Client()}

	err := service.replaceVerifiedBinary(context.Background(), executable, trusted.URL+"/asset", digest([]byte("foreign release")), "v1.0.1")
	if err == nil || !strings.Contains(err.Error(), "untrusted release redirect") {
		t.Fatalf("replaceVerifiedBinary() error = %v, want redirect refusal", err)
	}
}

func TestAttestationVerificationUsesExactReleaseWorkflow(t *testing.T) {
	t.Parallel()

	commands := &fakeCommands{responses: map[string]commandResponse{}}
	service := Service{Command: commands.Run}
	arguments := []string{
		"attestation", "verify", "/tmp/grat",
		"--repo", "phranck/grat",
		"--signer-workflow", "phranck/grat/.github/workflows/release.yml",
		"--source-ref", "refs/tags/v1.2.3",
		"--deny-self-hosted-runners",
	}
	commands.responses[commandKey("gh", arguments...)] = commandResponse{}

	if err := service.verifyArtifactAttestation(context.Background(), "/tmp/grat", "v1.2.3"); err != nil {
		t.Fatalf("verifyArtifactAttestation() error = %v", err)
	}
	if !commands.called(commandKey("gh", arguments...)) {
		t.Fatalf("commands = %#v, want constrained gh attestation verification", commands.calls)
	}
}

func releaseService(executable string, server *httptest.Server, goos string, goarch string) Service {
	return Service{
		Executable:   func() (string, error) { return executable, nil },
		EvalSymlinks: filepath.EvalSymlinks,
		Command: func(context.Context, string, ...string) ([]byte, error) {
			return nil, errors.New("not installed by Homebrew")
		},
		CurrentVersion:    func() string { return "v1.0.0" },
		ReleaseAPI:        server.URL,
		HTTPClient:        server.Client(),
		GOOS:              goos,
		GOARCH:            goarch,
		Rename:            os.Rename,
		VerifyAttestation: func(context.Context, string, string) error { return nil },
	}
}

func newReleaseServer(t *testing.T, goos string, goarch string, oldBinary []byte, newBinary []byte, corruptLatestChecksum bool) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
