package rizz

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	updateRepo          = "Gogoro/rizz"
	updateCheckInterval = 24 * time.Hour
	updateHTTPTimeout   = 5 * time.Second
	updateDownloadTimeout = 60 * time.Second
)

type updateCache struct {
	LatestVersion  string    `json:"latest_version"`
	CheckedAt      time.Time `json:"checked_at"`
	SkippedVersion string    `json:"skipped_version,omitempty"`
}

func updateCachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "rizz", "update-check.json")
}

func readUpdateCache() updateCache {
	path := updateCachePath()
	if path == "" {
		return updateCache{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return updateCache{}
	}
	var cache updateCache
	_ = json.Unmarshal(data, &cache)
	return cache
}

func writeUpdateCache(cache updateCache) {
	path := updateCachePath()
	if path == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}

func fetchLatestReleaseTag() (string, error) {
	url := "https://api.github.com/repos/" + updateRepo + "/releases/latest"
	client := &http.Client{Timeout: updateHTTPTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api: status %d", resp.StatusCode)
	}
	var body struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.TagName, nil
}

func parseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(strings.SplitN(parts[i], "-", 2)[0])
		out[i] = n
	}
	return out
}

func versionIsNewer(candidate, current string) bool {
	c := parseSemver(candidate)
	b := parseSemver(current)
	for i := 0; i < 3; i++ {
		if c[i] != b[i] {
			return c[i] > b[i]
		}
	}
	return false
}

// RefreshUpdateCacheBackground kicks off a non-blocking check against the
// GitHub releases API. The result is persisted to the cache for the *next*
// launch — we never block the user while waiting on the network.
func RefreshUpdateCacheBackground() {
	if Version == "dev" {
		return
	}
	cache := readUpdateCache()
	if time.Since(cache.CheckedAt) < updateCheckInterval {
		return
	}
	go func() {
		latest, err := fetchLatestReleaseTag()
		if err != nil {
			return
		}
		cache.LatestVersion = latest
		cache.CheckedAt = time.Now()
		writeUpdateCache(cache)
	}()
}

// PendingUpdateVersion returns the cached latest version if it is newer than
// the running version AND the user hasn't explicitly skipped it. Returns ""
// when there's nothing to prompt about.
func PendingUpdateVersion() string {
	if Version == "dev" {
		return ""
	}
	cache := readUpdateCache()
	if cache.LatestVersion == "" {
		return ""
	}
	if !versionIsNewer(cache.LatestVersion, Version) {
		return ""
	}
	if cache.SkippedVersion == cache.LatestVersion {
		return ""
	}
	return cache.LatestVersion
}

// RememberSkippedVersion persists that the user said "no" to this version so
// we don't nag them again until a new version comes out.
func RememberSkippedVersion(version string) {
	cache := readUpdateCache()
	cache.SkippedVersion = version
	writeUpdateCache(cache)
}

// SelfUpdate downloads the latest release tarball for the current OS/arch,
// extracts the rizz binary, and atomically replaces the running executable.
func SelfUpdate() error {
	if Version == "dev" {
		return fmt.Errorf("built from source; run: go install github.com/%s@latest", updateRepo)
	}
	if runtime.GOOS == "windows" {
		return fmt.Errorf("self-update not supported on windows; download from https://github.com/%s/releases", updateRepo)
	}

	latest, err := fetchLatestReleaseTag()
	if err != nil {
		return fmt.Errorf("fetch latest version: %w", err)
	}
	if !versionIsNewer(latest, Version) {
		fmt.Printf("already on the latest version (%s)\n", Version)
		return nil
	}

	url := releaseAssetURL(latest)
	fmt.Printf("downloading %s\n", url)
	binary, err := downloadReleaseBinary(url)
	if err != nil {
		return err
	}

	if err := replaceRunningBinary(binary); err != nil {
		return err
	}

	fmt.Printf("updated rizz to %s\n", latest)
	return nil
}

func releaseAssetURL(tag string) string {
	osName := runtime.GOOS
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}
	version := strings.TrimPrefix(tag, "v")
	asset := fmt.Sprintf("rizz_%s_%s_%s.tar.gz", version, osName, archName)
	return fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", updateRepo, tag, asset)
}

func downloadReleaseBinary(url string) ([]byte, error) {
	client := &http.Client{Timeout: updateDownloadTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: status %d", resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gunzip: %w", err)
	}
	defer gz.Close()

	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("untar: %w", err)
		}
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == "rizz" {
			return io.ReadAll(reader)
		}
	}
	return nil, fmt.Errorf("rizz binary not found in %s", url)
}

func replaceRunningBinary(newBinary []byte) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate current binary: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(exePath)
	if err == nil {
		exePath = resolved
	}

	tempPath := exePath + ".new"
	if err := os.WriteFile(tempPath, newBinary, 0755); err != nil {
		return fmt.Errorf("write new binary (is %s writable?): %w", filepath.Dir(exePath), err)
	}
	if err := os.Rename(tempPath, exePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("replace binary: %w (try re-running the install script)", err)
	}
	return nil
}
