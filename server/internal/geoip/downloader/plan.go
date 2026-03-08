package downloader

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/geoip"
)

const (
	dbipLatestArchiveName  = "dbip-country-lite.mmdb.gz"
	maxMindCountryDownload = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=%s&suffix=tar.gz"
)

var dbipMonthlyArchivePattern = regexp.MustCompile(`^dbip-country-lite-\d{4}-\d{2}\.mmdb\.gz$`)

func (d *Downloader) BuildDownloadPlan() (DownloadPlan, error) {
	switch d.ConfiguredSource() {
	case geoip.SourceDBIP:
		return d.buildDBIPPlan()
	case geoip.SourceDownloadURL:
		return DownloadPlan{
			Source:        geoip.SourceDownloadURL,
			CandidateURLs: []string{d.downloadURL},
		}, nil
	case geoip.SourceMaxMind:
		return DownloadPlan{
			Source:        geoip.SourceMaxMind,
			CandidateURLs: []string{fmt.Sprintf(maxMindCountryDownload, url.QueryEscape(d.maxMindLicenseKey))},
		}, nil
	default:
		return DownloadPlan{}, errors.New("GeoIP download source is not configured")
	}
}

func (d *Downloader) buildDBIPPlan() (DownloadPlan, error) {
	candidateURLs, err := buildDBIPCandidateURLs(d.downloadURL, time.Now().UTC())
	if err != nil {
		return DownloadPlan{}, err
	}

	return DownloadPlan{
		Source:        geoip.SourceDBIP,
		CandidateURLs: candidateURLs,
	}, nil
}

func buildDBIPCandidateURLs(downloadURL string, now time.Time) ([]string, error) {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("parse DB-IP download URL: %w", err)
	}

	archiveName := path.Base(parsedURL.Path)
	if archiveName == "." || archiveName == "/" || archiveName == "" {
		return nil, errors.New("DB-IP download URL must include a file name")
	}

	if dbipMonthlyArchivePattern.MatchString(archiveName) {
		return []string{downloadURL}, nil
	}

	if !strings.EqualFold(archiveName, dbipLatestArchiveName) {
		return nil, fmt.Errorf("unsupported DB-IP archive name %q", archiveName)
	}

	basePath := strings.TrimSuffix(parsedURL.Path, archiveName)
	candidateURLs := make([]string, 0, 4)
	for monthOffset := 0; monthOffset < 3; monthOffset++ {
		month := now.AddDate(0, -monthOffset, 0).Format("2006-01")
		candidateURL := *parsedURL
		candidateURL.Path = basePath + "dbip-country-lite-" + month + ".mmdb.gz"
		candidateURL.RawQuery = ""
		candidateURL.Fragment = ""
		candidateURLs = append(candidateURLs, candidateURL.String())
	}

	return appendUniqueURLs(candidateURLs, downloadURL), nil
}

func appendUniqueURLs(urls []string, extraURL string) []string {
	if extraURL == "" {
		return urls
	}

	for _, existingURL := range urls {
		if existingURL == extraURL {
			return urls
		}
	}

	return append(urls, extraURL)
}

func isDBIPDownloadURL(downloadURL string) bool {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return false
	}

	archiveName := strings.ToLower(path.Base(parsedURL.Path))
	return archiveName == dbipLatestArchiveName || dbipMonthlyArchivePattern.MatchString(archiveName)
}
