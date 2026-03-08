package downloader

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
)

func (d *Downloader) installDownloadedFile(downloadPath, downloadURL string) error {
	switch {
	case isTarGzDownload(downloadURL):
		return d.installDatabaseFromTarGz(downloadPath)
	case isGzipDownload(downloadURL):
		return d.installDatabaseFromGzip(downloadPath)
	default:
		return d.replaceDatabaseFile(downloadPath)
	}
}

func (d *Downloader) installDatabaseFromTarGz(archivePath string) error {
	archiveFile, err := os.Open(archivePath) // #nosec G304 -- path comes from our staged download file
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return fmt.Errorf("open GeoIP gzip archive: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		switch {
		case errors.Is(err, io.EOF):
			return errors.New("GeoIP archive does not contain an .mmdb file")
		case err != nil:
			return fmt.Errorf("read GeoIP tar archive: %w", err)
		case header == nil:
			continue
		case !header.FileInfo().Mode().IsRegular():
			continue
		case !strings.HasSuffix(strings.ToLower(header.Name), ".mmdb"):
			continue
		}

		return d.installDatabaseFromReader(tarReader)
	}
}

func (d *Downloader) installDatabaseFromGzip(archivePath string) error {
	archiveFile, err := os.Open(archivePath) // #nosec G304 -- path comes from our staged download file
	if err != nil {
		return fmt.Errorf("open GeoIP archive: %w", err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return fmt.Errorf("open GeoIP gzip archive: %w", err)
	}
	defer gzipReader.Close()

	return d.installDatabaseFromReader(gzipReader)
}

func (d *Downloader) installDatabaseFromReader(source io.Reader) error {
	stagedPath, err := d.writeStagedDatabase(source)
	if err != nil {
		return err
	}

	removeStagedFile := true
	defer func() {
		if removeStagedFile {
			_ = os.Remove(stagedPath)
		}
	}()

	if err := d.replaceDatabaseFile(stagedPath); err != nil {
		return err
	}

	removeStagedFile = false
	return nil
}

func (d *Downloader) writeStagedDatabase(source io.Reader) (string, error) {
	databaseFile, err := d.createTempFile("geoip-database-*")
	if err != nil {
		return "", err
	}

	databasePath := databaseFile.Name()
	if _, err := io.Copy(databaseFile, source); err != nil {
		_ = databaseFile.Close()
		_ = os.Remove(databasePath)
		return "", fmt.Errorf("write extracted GeoIP database: %w", err)
	}
	if err := databaseFile.Close(); err != nil {
		_ = os.Remove(databasePath)
		return "", fmt.Errorf("close extracted GeoIP database: %w", err)
	}

	return databasePath, nil
}

func isTarGzDownload(downloadURL string) bool {
	downloadTarget := downloadTarget(downloadURL)
	return downloadTarget == "tar.gz" ||
		downloadTarget == "tgz" ||
		strings.HasSuffix(downloadTarget, ".tar.gz") ||
		strings.HasSuffix(downloadTarget, ".tgz")
}

func isGzipDownload(downloadURL string) bool {
	if isTarGzDownload(downloadURL) {
		return false
	}

	downloadTarget := downloadTarget(downloadURL)
	return strings.HasSuffix(downloadTarget, ".gz") || downloadTarget == "gz" || downloadTarget == "gzip"
}

func downloadTarget(downloadURL string) string {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return strings.ToLower(downloadURL)
	}

	if suffix := strings.ToLower(parsedURL.Query().Get("suffix")); suffix != "" {
		return suffix
	}

	return strings.ToLower(path.Base(parsedURL.Path))
}
