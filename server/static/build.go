package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	if err := os.MkdirAll("dist", 0o755); err != nil {
		fmt.Printf("Failed to create dist directory: %v\n", err)
		os.Exit(1)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"tracker.ts"},
		Outfile:           "dist/tracker.js",
		Bundle:            true,
		Write:             true,
		Target:            api.ES2020,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		Platform:          api.PlatformBrowser,
		Format:            api.FormatIIFE,
		LegalComments:     api.LegalCommentsNone,
		LogLevel:          api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Build failed with %d errors\n", len(result.Errors))
		os.Exit(1)
	}

	if err := writeSizeBadge("dist/tracker.js", "dist/tracker-size.svg"); err != nil {
		fmt.Printf("Failed to write size badge: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Build successful!")
}

func writeSizeBadge(sourcePath, outputPath string) error {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	gzipSize, err := gzipLen(data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("tracker.js %s | gzip %s", formatBytes(len(data)), formatBytes(gzipSize))
	badge := renderBadge(text)
	return os.WriteFile(outputPath, []byte(badge), 0o644)
}

func gzipLen(data []byte) (int, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		_ = writer.Close()
		return 0, err
	}
	if err := writer.Close(); err != nil {
		return 0, err
	}
	return buf.Len(), nil
}

func formatBytes(size int) string {
	const unit = 1024.0
	value := float64(size)
	if value < unit {
		return fmt.Sprintf("%d B", size)
	}
	value = value / unit
	if value < unit {
		return fmt.Sprintf("%.1f KB", value)
	}
	value = value / unit
	return fmt.Sprintf("%.1f MB", value)
}

func renderBadge(text string) string {
	width := 40 + len(text)*7
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="26" viewBox="0 0 %d 26" role="img" aria-label="%s">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="#0b1220"/>
      <stop offset="1" stop-color="#111827"/>
    </linearGradient>
    <linearGradient id="stroke" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#22c55e"/>
      <stop offset="1" stop-color="#38bdf8"/>
    </linearGradient>
  </defs>
  <rect width="%d" height="26" rx="8" fill="url(#bg)"/>
  <rect x="0.5" y="0.5" width="%d" height="25" rx="7.5" fill="none" stroke="url(#stroke)" stroke-opacity="0.7"/>
  <text x="16" y="17" fill="#f8fafc" font-family="SFMono-Regular, Menlo, Consolas, monospace" font-size="12" letter-spacing="0.2">%s</text>
</svg>
`, width, width, text, width, width-1, text)
}
