package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/geoip/downloader"
	"github.com/stretchr/testify/require"
)

func TestGeoIPService_InvalidPath(t *testing.T) {
	t.Parallel()

	svc := NewService(Config{DBPath: "does-not-exist.mmdb"})
	require.NoError(t, svc.SetEnabled(true))

	got, err := svc.ResolveCountry("8.8.8.8")
	require.Equal(t, UnknownCountry, got)
	require.ErrorIs(t, err, ErrNoDBReader)
}

func TestGeoIPService_EnsureAvailable_MissingSource(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc := NewService(Config{DBPath: dbPath})

	require.NoError(t, svc.SetEnabled(true))
	err := svc.EnsureAvailable(context.Background())
	require.Error(t, err)

	status := svc.Status()
	require.Equal(t, StateMissing, status.State)
}

func TestGeoIPService_EnsureAvailable_DownloadFailure(t *testing.T) {
	t.Parallel()

	server := newIPv4TestServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc := NewService(Config{
		DBPath:      dbPath,
		DownloadURL: server.URL + "/GeoLite2-Country.tar.gz",
	})

	require.NoError(t, svc.SetEnabled(true))
	err := svc.EnsureAvailable(context.Background())
	require.Error(t, err)

	status := svc.Status()
	require.Equal(t, StateError, status.State)
	require.NotEmpty(t, status.LastError)
}

func TestGeoIPService_EnsureAvailable_DownloadInvalidArchive(t *testing.T) {
	t.Parallel()

	archive := buildTarGz(t, "GeoLite2-Country.mmdb", []byte("not-a-real-mmdb"))
	server := newIPv4TestServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "GeoLite2-Country.mmdb")
	svc := NewService(Config{
		DBPath:      dbPath,
		DownloadURL: server.URL + "/GeoLite2-Country.tar.gz",
	})

	require.NoError(t, svc.SetEnabled(true))
	err := svc.EnsureAvailable(context.Background())
	require.Error(t, err)

	_, statErr := os.Stat(dbPath)
	require.NoError(t, statErr)
}

func TestGeoIPService_EnsureAvailable_DoesNotRefreshLoadedReader(t *testing.T) {
	t.Parallel()

	now := time.Now()
	lookup := &fakeGeoIPLookup{
		hasReader: true,
		updatedAt: &now,
	}
	downloader := &fakeGeoIPDownloader{
		configured: true,
		source:     SourceMaxMind,
		plan: downloader.DownloadPlan{
			CandidateURLs: []string{"https://example.com/db.tar.gz"},
			Source:        SourceMaxMind,
		},
	}

	svc := NewService(Config{DBPath: "/tmp/test.mmdb"})
	svc.lookup = lookup
	svc.downloader = downloader
	require.NoError(t, svc.SetEnabled(true))

	err := svc.EnsureAvailable(context.Background())
	require.NoError(t, err)
	require.Zero(t, lookup.loadCalls)
	require.Zero(t, downloader.downloadCalls)

	status := svc.Status()
	require.Equal(t, StateReady, status.State)
	require.Equal(t, SourceFile, status.Source)
}

func TestGeoIPService_Refresh_ForcesReload(t *testing.T) {
	t.Parallel()

	now := time.Now()
	lookup := &fakeGeoIPLookup{
		hasReader: true,
		updatedAt: &now,
	}
	downloader := &fakeGeoIPDownloader{
		configured: true,
		source:     SourceDBIP,
		plan: downloader.DownloadPlan{
			CandidateURLs: []string{"https://example.com/db.tar.gz"},
			Source:        SourceDBIP,
		},
	}

	svc := NewService(Config{DBPath: "/tmp/test.mmdb"})
	svc.lookup = lookup
	svc.downloader = downloader

	err := svc.Refresh(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, downloader.planCalls)
	require.Equal(t, 1, downloader.downloadCalls)
	require.Equal(t, 1, lookup.loadCalls)

	status := svc.Status()
	require.Equal(t, StateReady, status.State)
	require.Equal(t, SourceDBIP, status.Source)
}

func TestGeoIPService_Refresh_PreservesErrorCause(t *testing.T) {
	t.Parallel()

	sentinel := context.Canceled
	downloader := &fakeGeoIPDownloader{
		configured: true,
		source:     SourceMaxMind,
		plan: downloader.DownloadPlan{
			CandidateURLs: []string{"https://example.com/db.tar.gz"},
			Source:        SourceMaxMind,
		},
		downloadErr: sentinel,
	}

	svc := NewService(Config{DBPath: "/tmp/test.mmdb"})
	svc.lookup = &fakeGeoIPLookup{}
	svc.downloader = downloader

	err := svc.Refresh(context.Background())
	require.ErrorIs(t, err, sentinel)

	status := svc.Status()
	require.Equal(t, StateError, status.State)
	require.Contains(t, status.LastError, sentinel.Error())
}

func TestGeoIPService_DisablePreservesCloseFailure(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("close failed")
	lookup := &fakeGeoIPLookup{hasReader: true, closeErr: sentinel}
	svc := NewService(Config{DBPath: "/tmp/test.mmdb"})
	svc.lookup = lookup
	require.NoError(t, svc.SetEnabled(true))

	err := svc.SetEnabled(false)
	require.ErrorIs(t, err, sentinel)
	require.Equal(t, StateError, svc.Status().State)
	require.Contains(t, svc.Status().LastError, sentinel.Error())
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
	require.NoError(t, tarWriter.WriteHeader(header))
	_, err := tarWriter.Write(contents)
	require.NoError(t, err)
	require.NoError(t, tarWriter.Close())
	require.NoError(t, gzipWriter.Close())
	return buf.Bytes()
}

func newIPv4TestServer(handler http.Handler) *httptest.Server {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Errorf("listen on 127.0.0.1:0: %w", err))
	}
	server := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second},
	}
	server.Start()
	return server
}

type fakeGeoIPLookup struct {
	hasReader      bool
	fileExists     bool
	updatedAt      *time.Time
	loadErr        error
	loadCalls      int
	closeErr       error
	closeCalls     int
	listCountries  []ListedCountry
	listErr        error
	resolveCountry Country
	resolveErr     error
}

func (f *fakeGeoIPLookup) HasReader() bool {
	return f.hasReader
}

func (f *fakeGeoIPLookup) FileExists() bool {
	return f.fileExists
}

func (f *fakeGeoIPLookup) UpdatedAt() *time.Time {
	return f.updatedAt
}

func (f *fakeGeoIPLookup) Load() error {
	f.loadCalls++
	if f.loadErr != nil {
		return f.loadErr
	}
	f.hasReader = true
	return nil
}

func (f *fakeGeoIPLookup) ListCountries(string) ([]ListedCountry, error) {
	return f.listCountries, f.listErr
}

func (f *fakeGeoIPLookup) ResolveCountry(string) (Country, error) {
	return f.resolveCountry, f.resolveErr
}

func (f *fakeGeoIPLookup) Close() error {
	f.closeCalls++
	f.hasReader = false
	return f.closeErr
}

type fakeGeoIPDownloader struct {
	configured    bool
	source        Source
	plan          downloader.DownloadPlan
	planErr       error
	planCalls     int
	downloadErr   error
	downloadCalls int
}

func (f *fakeGeoIPDownloader) HasDownloadSource() bool {
	return f.configured
}

func (f *fakeGeoIPDownloader) ConfiguredSource() Source {
	return f.source
}

func (f *fakeGeoIPDownloader) BuildDownloadPlan() (downloader.DownloadPlan, error) {
	f.planCalls++
	if f.planErr != nil {
		return downloader.DownloadPlan{}, f.planErr
	}
	return f.plan, nil
}

func (f *fakeGeoIPDownloader) Download(context.Context, downloader.DownloadPlan) error {
	f.downloadCalls++
	return f.downloadErr
}

var _ geoIPLookup = (*fakeGeoIPLookup)(nil)
var _ geoIPDownloader = (*fakeGeoIPDownloader)(nil)
