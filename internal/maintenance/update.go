package maintenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxReleaseDocumentBytes int64 = 1 << 20
	defaultMaxReleaseAssetBytes    int64 = 128 << 20
	releaseRepository                    = "phranck/grat"
	releaseSignerWorkflow                = "phranck/grat/.github/workflows/release.yml"
	releaseAPIPathPrefix                 = "/repos/phranck/grat/releases/"
	releaseAssetPathPrefix               = "/phranck/grat/releases/download/"
)

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Update updates the current grat installation only when its owner can be
// determined safely.
func (service Service) Update(ctx context.Context) (Result, error) {
	executable, err := service.executable()
	if err != nil {
		return Result{}, err
	}
	if owned, err := service.homebrewOwns(ctx, executable); err != nil {
		return Result{}, err
	} else if owned {
		if _, err := service.command(ctx, "brew", "upgrade", HomebrewFormula); err != nil {
			return Result{}, fmt.Errorf("update with Homebrew: %w", err)
		}
		return Result{Message: "Updated grat with Homebrew."}, nil
	}
	if module, buildVersion, ok := service.buildInfo(); ok && module == ModulePath && buildVersion != "" && buildVersion != "(devel)" {
		return Result{Message: "Run: go install github.com/phranck/grat/cmd/grat@latest"}, nil
	}
	return service.updateDirectRelease(ctx, executable)
}

func (service Service) homebrewOwns(ctx context.Context, executable string) (bool, error) {
	if _, err := service.command(ctx, "brew", "list", "--versions", HomebrewFormula); err != nil {
		return false, nil
	}
	prefixOutput, err := service.command(ctx, "brew", "--prefix", HomebrewFormula)
	if err != nil {
		return false, nil
	}
	prefix := strings.TrimSpace(string(prefixOutput))
	if prefix == "" {
		return false, nil
	}
	current, err := service.evalSymlinks(executable)
	if err != nil {
		return false, fmt.Errorf("resolve current executable: %w", err)
	}
	candidate, err := service.evalSymlinks(filepath.Join(prefix, "bin", "grat"))
	if err != nil {
		return false, nil
	}
	return current == candidate, nil
}

func (service Service) updateDirectRelease(ctx context.Context, executable string) (Result, error) {
	currentVersion := service.currentVersion()
	if currentVersion == "" {
		return Result{}, errors.New("cannot determine the installed grat version")
	}
	if err := service.verifyDirectRelease(ctx, executable, currentVersion); err != nil {
		return Result{}, err
	}

	latest, err := service.release(ctx, "latest")
	if err != nil {
		return Result{}, fmt.Errorf("load latest grat release: %w", err)
	}
	if latest.TagName == "" {
		return Result{}, errors.New("latest grat release has no tag")
	}
	if latest.TagName == currentVersion {
		return Result{Message: "grat is already up to date (" + currentVersion + ")."}, nil
	}
	assetName := service.assetName(latest.TagName)
	expectedDigest, err := service.releaseChecksum(ctx, latest, assetName)
	if err != nil {
		return Result{}, fmt.Errorf("load checksum for %s: %w", assetName, err)
	}
	binaryAsset, ok := findAsset(latest.Assets, assetName)
	if !ok {
		return Result{}, fmt.Errorf("latest grat release has no asset %s", assetName)
	}
	if err := service.replaceVerifiedBinary(ctx, executable, binaryAsset.BrowserDownloadURL, expectedDigest, latest.TagName); err != nil {
		return Result{}, err
	}
	return Result{Message: "Updated grat to " + latest.TagName + "."}, nil
}

func (service Service) verifyDirectRelease(ctx context.Context, executable string, currentVersion string) error {
	currentRelease, err := service.release(ctx, "tags/"+url.PathEscape(currentVersion))
	if err != nil {
		return fmt.Errorf("verify installed release %s: %w", currentVersion, err)
	}
	currentAssetName := service.assetName(currentVersion)
	currentExpected, err := service.releaseChecksum(ctx, currentRelease, currentAssetName)
	if err != nil {
		return fmt.Errorf("verify installed release %s: %w", currentVersion, err)
	}
	currentDigest, err := fileDigest(executable)
	if err != nil {
		return fmt.Errorf("hash current executable: %w", err)
	}
	if !strings.EqualFold(currentDigest, currentExpected) {
		return errors.New("current executable does not match the verified GitHub release checksum")
	}
	if err := service.verifyArtifactAttestation(ctx, executable, currentVersion); err != nil {
		return fmt.Errorf("verify current release provenance: %w", err)
	}
	return nil
}

func (service Service) release(ctx context.Context, suffix string) (release, error) {
	endpoint, err := service.endpoint("/repos/phranck/grat/releases/" + suffix)
	if err != nil {
		return release{}, err
	}
	data, err := service.get(ctx, endpoint, false)
	if err != nil {
		return release{}, err
	}
	var value release
	if err := json.Unmarshal(data, &value); err != nil {
		return release{}, fmt.Errorf("parse release metadata: %w", err)
	}
	return value, nil
}

func (service Service) releaseChecksum(ctx context.Context, value release, assetName string) (string, error) {
	checksums, ok := findAsset(value.Assets, "checksums.txt")
	if !ok {
		return "", errors.New("release has no checksums.txt")
	}
	endpoint, err := service.assetURL(checksums.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	data, err := service.get(ctx, endpoint, true)
	if err != nil {
		return "", err
	}
	return checksumForAsset(string(data), assetName)
}

func (service Service) replaceVerifiedBinary(ctx context.Context, executable string, assetURL string, expectedDigest string, tag string) error {
	endpoint, err := service.assetURL(assetURL)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create release download request: %w", err)
	}
	client, err := service.releaseHTTPClient(true)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("download release asset: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download release asset: unexpected HTTP status %s", response.Status)
	}
	limit := service.maxReleaseAssetBytes()
	if err := validateResponseLength(response, limit, "release asset"); err != nil {
		return err
	}
	info, err := os.Stat(executable)
	if err != nil {
		return fmt.Errorf("inspect current executable: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(executable), ".grat-update-*")
	if err != nil {
		return fmt.Errorf("create temporary update: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		return errors.Join(fmt.Errorf("set temporary update permissions: %w", err), temporary.Close())
	}
	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(temporary, hasher), io.LimitReader(response.Body, limit+1))
	if err != nil {
		return errors.Join(fmt.Errorf("write temporary update: %w", err), temporary.Close())
	}
	if written > limit {
		return errors.Join(fmt.Errorf("release asset exceeds maximum size of %d bytes", limit), temporary.Close())
	}
	if actual := hex.EncodeToString(hasher.Sum(nil)); !strings.EqualFold(actual, expectedDigest) {
		return errors.Join(errors.New("downloaded release checksum does not match checksums.txt"), temporary.Close())
	}
	if err := temporary.Sync(); err != nil {
		return errors.Join(fmt.Errorf("sync temporary update: %w", err), temporary.Close())
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary update: %w", err)
	}
	if err := service.verifyArtifactAttestation(ctx, temporaryPath, tag); err != nil {
		return fmt.Errorf("verify downloaded release provenance: %w", err)
	}
	if err := service.rename(temporaryPath, executable); err != nil {
		return fmt.Errorf("replace executable: %w", err)
	}
	return nil
}

func (service Service) endpoint(path string) (string, error) {
	baseURL, err := service.releaseAPIURL()
	if err != nil {
		return "", err
	}
	reference, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("parse release API path: %w", err)
	}
	endpoint := baseURL.ResolveReference(reference)
	if err := service.validateReleaseAPIURL(endpoint); err != nil {
		return "", err
	}
	return endpoint.String(), nil
}

func (service Service) assetURL(raw string) (string, error) {
	baseURL, err := service.releaseAPIURL()
	if err != nil {
		return "", err
	}
	reference, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse release asset URL: %w", err)
	}
	endpoint := baseURL.ResolveReference(reference)
	if err := service.validateReleaseAssetURL(endpoint, false); err != nil {
		return "", err
	}
	return endpoint.String(), nil
}

func (service Service) get(ctx context.Context, endpoint string, asset bool) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	client, err := service.releaseHTTPClient(asset)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected HTTP status %s", response.Status)
	}
	return readBoundedResponse(response, service.maxReleaseDocumentBytes(), "release document")
}

func (service Service) releaseAPIURL() (*url.URL, error) {
	base := service.ReleaseAPI
	if base == "" {
		base = defaultReleaseAPI
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse release API URL: %w", err)
	}
	if err := validateHTTPSURL(baseURL, "release API"); err != nil {
		return nil, err
	}
	return baseURL, nil
}

func (service Service) validateReleaseAPIURL(endpoint *url.URL) error {
	baseURL, err := service.releaseAPIURL()
	if err != nil {
		return err
	}
	if err := validateHTTPSURL(endpoint, "release API"); err != nil {
		return err
	}
	if endpoint.Scheme != baseURL.Scheme || endpoint.Host != baseURL.Host {
		return fmt.Errorf("untrusted release API origin %q", endpoint.Host)
	}
	if service.productionReleaseAPI(baseURL) && !strings.HasPrefix(endpoint.Path, releaseAPIPathPrefix) {
		return fmt.Errorf("untrusted release API path %q", endpoint.Path)
	}
	return nil
}

func (service Service) validateReleaseAssetURL(endpoint *url.URL, redirect bool) error {
	baseURL, err := service.releaseAPIURL()
	if err != nil {
		return err
	}
	if err := validateHTTPSURL(endpoint, "release asset"); err != nil {
		return err
	}
	if !service.productionReleaseAPI(baseURL) {
		if endpoint.Scheme != baseURL.Scheme || endpoint.Host != baseURL.Host {
			return fmt.Errorf("untrusted release asset origin %q", endpoint.Host)
		}
		return nil
	}

	host := endpoint.Hostname()
	if endpoint.Port() != "" {
		return fmt.Errorf("untrusted release asset origin %q", endpoint.Host)
	}
	if host == "github.com" && strings.HasPrefix(endpoint.Path, releaseAssetPathPrefix) {
		return nil
	}
	if redirect && (host == "release-assets.githubusercontent.com" || host == "objects.githubusercontent.com") {
		return nil
	}
	return fmt.Errorf("untrusted release asset origin %q", endpoint.Host)
}

func validateHTTPSURL(endpoint *url.URL, description string) error {
	if endpoint == nil || endpoint.Scheme != "https" || endpoint.Host == "" || endpoint.User != nil || endpoint.Fragment != "" {
		return fmt.Errorf("%s URL must use an absolute credential-free HTTPS origin", description)
	}
	return nil
}

func (service Service) productionReleaseAPI(endpoint *url.URL) bool {
	return endpoint.Scheme == "https" && endpoint.Hostname() == "api.github.com" && endpoint.Port() == ""
}

func (service Service) releaseHTTPClient(asset bool) (*http.Client, error) {
	client := *service.httpClient()
	previousCheck := client.CheckRedirect
	client.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		var err error
		if asset {
			err = service.validateReleaseAssetURL(request.URL, true)
		} else {
			err = service.validateReleaseAPIURL(request.URL)
		}
		if err != nil {
			return fmt.Errorf("untrusted release redirect: %w", err)
		}
		if previousCheck != nil {
			return previousCheck(request, via)
		}
		return nil
	}
	return &client, nil
}

func (service Service) verifyArtifactAttestation(ctx context.Context, path string, tag string) error {
	if service.VerifyAttestation != nil {
		return service.VerifyAttestation(ctx, path, tag)
	}
	_, err := service.command(
		ctx,
		"gh",
		"attestation", "verify", path,
		"--repo", releaseRepository,
		"--signer-workflow", releaseSignerWorkflow,
		"--source-ref", "refs/tags/"+tag,
		"--deny-self-hosted-runners",
	)
	if err != nil {
		return fmt.Errorf("verify GitHub artifact attestation for %s: %w", tag, err)
	}
	return nil
}

func readBoundedResponse(response *http.Response, limit int64, description string) ([]byte, error) {
	if err := validateResponseLength(response, limit, description); err != nil {
		return nil, err
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%s exceeds maximum size of %d bytes", description, limit)
	}
	return data, nil
}

func validateResponseLength(response *http.Response, limit int64, description string) error {
	if response.ContentLength > limit {
		return fmt.Errorf("%s exceeds maximum size of %d bytes", description, limit)
	}
	return nil
}

func (service Service) assetName(tag string) string {
	return fmt.Sprintf("grat_%s_%s_%s", tag, service.goos(), service.goarch())
}

func (service Service) executable() (string, error) {
	if service.Executable == nil {
		return os.Executable()
	}
	return service.Executable()
}

func (service Service) evalSymlinks(path string) (string, error) {
	if service.EvalSymlinks == nil {
		return filepath.EvalSymlinks(path)
	}
	return service.EvalSymlinks(path)
}

func (service Service) command(ctx context.Context, name string, arguments ...string) ([]byte, error) {
	if service.Command == nil {
		return runCommand(ctx, name, arguments...)
	}
	return service.Command(ctx, name, arguments...)
}

func (service Service) buildInfo() (string, string, bool) {
	if service.BuildInfo == nil {
		return runningBuildInfo()
	}
	return service.BuildInfo()
}

func (service Service) currentVersion() string {
	if service.CurrentVersion == nil {
		return ""
	}
	return service.CurrentVersion()
}

func (service Service) goos() string {
	return service.GOOS
}

func (service Service) goarch() string {
	return service.GOARCH
}

func (service Service) rename(oldPath string, newPath string) error {
	if service.Rename == nil {
		return os.Rename(oldPath, newPath)
	}
	return service.Rename(oldPath, newPath)
}

func (service Service) httpClient() *http.Client {
	if service.HTTPClient != nil {
		return service.HTTPClient
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (service Service) maxReleaseDocumentBytes() int64 {
	if service.MaxReleaseDocumentBytes > 0 {
		return service.MaxReleaseDocumentBytes
	}
	return defaultMaxReleaseDocumentBytes
}

func (service Service) maxReleaseAssetBytes() int64 {
	if service.MaxReleaseAssetBytes > 0 {
		return service.MaxReleaseAssetBytes
	}
	return defaultMaxReleaseAssetBytes
}

func findAsset(assets []asset, name string) (asset, bool) {
	for _, asset := range assets {
		if asset.Name == name {
			return asset, true
		}
	}
	return asset{}, false
}

func checksumForAsset(document string, assetName string) (string, error) {
	for _, line := range strings.Split(document, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || strings.TrimPrefix(fields[1], "*") != assetName {
			continue
		}
		if len(fields[0]) != sha256.Size*2 {
			return "", fmt.Errorf("invalid checksum for %s", assetName)
		}
		if _, err := hex.DecodeString(fields[0]); err != nil {
			return "", fmt.Errorf("invalid checksum for %s", assetName)
		}
		return fields[0], nil
	}
	return "", fmt.Errorf("checksums.txt has no checksum for %s", assetName)
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
