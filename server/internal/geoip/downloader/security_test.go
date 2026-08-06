package downloader

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyWithLimitRejectsOversizedPayload(t *testing.T) {
	var destination bytes.Buffer

	err := copyWithLimit(&destination, strings.NewReader("123456789"), 8)

	require.ErrorIs(t, err, errPayloadTooLarge)
	require.Len(t, destination.Bytes(), 9)
}

func TestDownloadURLForLogRedactsCredentialsQueryAndFragment(t *testing.T) {
	logged := downloadURLForLog("https://user:secret@download.example/db.mmdb?license_key=private#fragment")

	require.Equal(t, "https://download.example", logged)
}
