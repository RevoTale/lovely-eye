package services

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestGeoIPService_InvalidPath(t *testing.T) {
	t.Parallel()

	svc, err := NewGeoIPService(GeoIPConfig{DBPath: "does-not-exist.mmdb"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	svc.SetEnabled(true)

	if got,err := svc.ResolveCountry("8.8.8.8"); got.Name != "Unknown"  {
		t.Fatalf("expected Unknown for public IP with missing DB, got %q %s", got,err.Error())
		require.NoError(t,err)
	}
}

func TestGeoIPService_EnsureAvailable_MissingSource(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc, err := NewGeoIPService(GeoIPConfig{DBPath: dbPath})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	svc.SetEnabled(true)
	if err := svc.EnsureAvailable(context.Background()); err == nil {
		t.Fatalf("expected error when download source is missing")
	}

	status := svc.Status()
	if status.State != geoIPStateMissing {
		t.Fatalf("expected state %q, got %q", geoIPStateMissing, status.State)
	}
}

func TestGeoIPService_EnsureAvailable_DownloadFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc, err := NewGeoIPService(GeoIPConfig{
		DBPath:      dbPath,
		DownloadURL: server.URL + "/GeoLite2-Country.tar.gz",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	svc.SetEnabled(true)
	if err := svc.EnsureAvailable(context.Background()); err == nil {
		t.Fatalf("expected download error")
	}

	status := svc.Status()
	if status.State != geoIPStateError {
		t.Fatalf("expected state %q, got %q", geoIPStateError, status.State)
	}
	if status.LastError == "" {
		t.Fatalf("expected last error to be set")
	}
}

func TestGeoIPService_EnsureAvailable_DownloadInvalidArchive(t *testing.T) {
	t.Parallel()

	archive := buildTarGz(t, "GeoLite2-Country.mmdb", []byte("not-a-real-mmdb"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc, err := NewGeoIPService(GeoIPConfig{
		DBPath:      dbPath,
		DownloadURL: server.URL + "/GeoLite2-Country.tar.gz",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	svc.SetEnabled(true)
	if err := svc.EnsureAvailable(context.Background()); err == nil {
		t.Fatalf("expected load error for invalid database")
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file to be written, got %v", err)
	}
}

func buildTarGz(t *testing.T, name string, contents []byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	header := &tar.Header{
		Name: name,
		Mode: 0o644,
		Size: int64(len(contents)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("write tar header: %v", err)
	}
	if _, err := tarWriter.Write(contents); err != nil {
		t.Fatalf("write tar data: %v", err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}
	return buf.Bytes()
}
