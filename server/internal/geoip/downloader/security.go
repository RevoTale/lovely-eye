package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/url"
)

const (
	maxDownloadedBytes = 64 << 20
	maxExtractedBytes  = 128 << 20
)

var errPayloadTooLarge = errors.New("GeoIP payload is too large")

func copyWithLimit(destination io.Writer, source io.Reader, maximum int64) error {
	written, err := io.Copy(destination, io.LimitReader(source, maximum+1))
	if err != nil {
		return fmt.Errorf("copy GeoIP payload: %w", err)
	}
	if written > maximum {
		return errPayloadTooLarge
	}
	return nil
}

func downloadURLForLog(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "<invalid URL>"
	}
	parsed.User = nil
	parsed.Path = ""
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.ForceQuery = false
	parsed.Fragment = ""
	return parsed.String()
}

func downloadRequestCause(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Err != nil {
		return urlErr.Err
	}
	return err
}
