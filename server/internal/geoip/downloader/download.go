package downloader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func (d *Downloader) downloadAndInstall(ctx context.Context, plan DownloadPlan) error {
	if len(plan.CandidateURLs) == 0 {
		return fmt.Errorf("GeoIP download plan for %s has no candidate URLs", plan.Source)
	}

	var lastErr error
	for _, candidateURL := range plan.CandidateURLs {
		if err := d.downloadCandidate(ctx, candidateURL); err != nil {
			lastErr = err
			slog.Warn("failed to download GeoIP database", "url", candidateURL, "error", err)
			continue
		}
		return nil
	}

	return lastErr
}

func (d *Downloader) downloadCandidate(ctx context.Context, downloadURL string) error {
	downloadPath, err := d.downloadToTempFile(ctx, downloadURL)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(downloadPath)
	}()

	if err := d.installDownloadedFile(downloadPath, downloadURL); err != nil {
		return fmt.Errorf("install GeoIP database from %q: %w", downloadURL, err)
	}

	return nil
}

func (d *Downloader) downloadToTempFile(ctx context.Context, downloadURL string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("build GeoIP download request: %w", err)
	}

	response, err := d.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("download GeoIP database from %q: %w", downloadURL, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download GeoIP database from %q: unexpected status %s", downloadURL, response.Status)
	}

	downloadFile, err := d.createTempFile("geoip-download-*")
	if err != nil {
		return "", err
	}

	downloadPath := downloadFile.Name()
	if _, err := io.Copy(downloadFile, response.Body); err != nil {
		_ = downloadFile.Close()
		_ = os.Remove(downloadPath)
		return "", fmt.Errorf("write downloaded GeoIP payload: %w", err)
	}
	if err := downloadFile.Close(); err != nil {
		_ = os.Remove(downloadPath)
		return "", fmt.Errorf("close downloaded GeoIP payload: %w", err)
	}

	return downloadPath, nil
}

func (d *Downloader) createTempFile(pattern string) (*os.File, error) {
	databaseDir := filepath.Dir(d.dbPath)
	if err := os.MkdirAll(databaseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create GeoIP database directory: %w", err)
	}

	file, err := os.CreateTemp(databaseDir, pattern)
	if err != nil {
		return nil, fmt.Errorf("create temporary GeoIP file: %w", err)
	}

	return file, nil
}

func (d *Downloader) replaceDatabaseFile(stagedPath string) error {
	if err := os.Rename(stagedPath, d.dbPath); err != nil {
		return fmt.Errorf("replace GeoIP database at %q: %w", d.dbPath, err)
	}
	return nil
}
