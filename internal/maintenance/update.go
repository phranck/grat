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
	if err := service.replaceVerifiedBinary(ctx, executable, binaryAsset.BrowserDownloadURL, expectedDigest); err != nil {
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
	return nil
}

func (service Service) release(ctx context.Context, suffix string) (release, error) {
	endpoint, err := service.endpoint("/repos/phranck/grat/releases/" + suffix)
	if err != nil {
		return release{}, err
	}
	data, err := service.get(ctx, endpoint)
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
	data, err := service.get(ctx, endpoint)
	if err != nil {
		return "", err
	}
	return checksumForAsset(string(data), assetName)
}

func (service Service) replaceVerifiedBinary(ctx context.Context, executable string, assetURL string, expectedDigest string) error {
	endpoint, err := service.assetURL(assetURL)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create release download request: %w", err)
	}
	response, err := service.httpClient().Do(request)
	if err != nil {
		return fmt.Errorf("download release asset: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download release asset: unexpected HTTP status %s", response.Status)
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
	if _, err := io.Copy(io.MultiWriter(temporary, hasher), response.Body); err != nil {
		return errors.Join(fmt.Errorf("write temporary update: %w", err), temporary.Close())
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
	if err := service.rename(temporaryPath, executable); err != nil {
		return fmt.Errorf("replace executable: %w", err)
	}
	return nil
}

func (service Service) endpoint(path string) (string, error) {
	base := service.ReleaseAPI
	if base == "" {
		base = defaultReleaseAPI
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse release API URL: %w", err)
	}
	reference, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("parse release API path: %w", err)
	}
	return baseURL.ResolveReference(reference).String(), nil
}

func (service Service) assetURL(raw string) (string, error) {
	return service.endpoint(raw)
}

func (service Service) get(ctx context.Context, endpoint string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := service.httpClient().Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected HTTP status %s", response.Status)
	}
	return io.ReadAll(io.LimitReader(response.Body, 8<<20))
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
