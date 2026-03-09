package models

import "testing"

func TestClientBrowserFromLegacyLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ClientBrowser
	}{
		{
			name:     "canonical chrome",
			input:    "Chrome",
			expected: ClientBrowserChrome,
		},
		{
			name:     "legacy chrome mobile",
			input:    "chrome mobile ios",
			expected: ClientBrowserChrome,
		},
		{
			name:     "unknown non-empty maps to other",
			input:    "Arc",
			expected: ClientBrowserOther,
		},
		{
			name:     "empty maps to unknown",
			input:    "",
			expected: ClientBrowserUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientBrowserFromLegacyLabel(tt.input); got != tt.expected {
				t.Fatalf("ClientBrowserFromLegacyLabel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestClientOSFromLegacyLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ClientOS
	}{
		{
			name:     "canonical ios",
			input:    "iOS",
			expected: ClientOSIOS,
		},
		{
			name:     "mac os x alias",
			input:    "Mac OS X",
			expected: ClientOSMacOS,
		},
		{
			name:     "unknown non-empty maps to other",
			input:    "Haiku",
			expected: ClientOSOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientOSFromLegacyLabel(tt.input); got != tt.expected {
				t.Fatalf("ClientOSFromLegacyLabel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestClientScreenSizeFromLegacyLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ClientScreenSize
	}{
		{
			name:     "canonical size",
			input:    "md",
			expected: ClientScreenSizeMD,
		},
		{
			name:     "legacy resolution",
			input:    "390x844",
			expected: ClientScreenSizeXS,
		},
		{
			name:     "watch resolution",
			input:    "280x280",
			expected: ClientScreenSizeWatch,
		},
		{
			name:     "invalid label",
			input:    "unknown-size",
			expected: ClientScreenSizeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClientScreenSizeFromLegacyLabel(tt.input); got != tt.expected {
				t.Fatalf("ClientScreenSizeFromLegacyLabel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseClientDeviceFilters(t *testing.T) {
	got := ParseClientDeviceFilters([]string{"mobile", "desktop", "mobile", "unknown"})
	want := []ClientDevice{ClientDeviceMobile, ClientDeviceDesktop}

	if len(got) != len(want) {
		t.Fatalf("ParseClientDeviceFilters() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseClientDeviceFilters()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
