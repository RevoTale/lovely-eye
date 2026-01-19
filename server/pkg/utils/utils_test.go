package utils

import (
	"errors"
	"testing"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError error
	}{
		{
			name:      "simple domain",
			input:     "example.com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "subdomain",
			input:     "blog.example.com",
			want:      "blog.example.com",
			wantError: nil,
		},
		{
			name:      "domain with https",
			input:     "https://example.com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with http",
			input:     "http://example.com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with www",
			input:     "www.example.com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with https and www",
			input:     "https://www.example.com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with path",
			input:     "example.com/path/to/page",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with trailing slash",
			input:     "example.com/",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "uppercase domain",
			input:     "EXAMPLE.COM",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "mixed case domain",
			input:     "Example.Com",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "domain with hyphens",
			input:     "my-site.example.com",
			want:      "my-site.example.com",
			wantError: nil,
		},
		{
			name:      "domain with whitespace",
			input:     "  example.com  ",
			want:      "example.com",
			wantError: nil,
		},
		{
			name:      "complex normalization",
			input:     "  HTTPS://WWW.Example-Site.COM/path  ",
			want:      "example-site.com",
			wantError: nil,
		},
		{
			name:      "empty domain",
			input:     "",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "only whitespace",
			input:     "   ",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "invalid domain - starts with hyphen",
			input:     "-example.com",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "invalid domain - ends with hyphen",
			input:     "example-.com",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "invalid domain - double dot",
			input:     "example..com",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "invalid domain - special characters",
			input:     "example$.com",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "invalid domain - spaces in middle",
			input:     "exam ple.com",
			want:      "",
			wantError: ErrInvalidDomain,
		},
		{
			name:      "single word domain",
			input:     "localhost",
			want:      "localhost",
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDomain(tt.input)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("ValidateDomain() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSiteName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError error
	}{
		{
			name:      "valid site name",
			input:     "My Website",
			want:      "My Website",
			wantError: nil,
		},
		{
			name:      "single character",
			input:     "A",
			want:      "A",
			wantError: nil,
		},
		{
			name:      "exactly 100 characters",
			input:     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			want:      "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			wantError: nil,
		},
		{
			name:      "with whitespace at edges",
			input:     "  My Website  ",
			want:      "My Website",
			wantError: nil,
		},
		{
			name:      "with special characters",
			input:     "My Website - Analytics & Stats!",
			want:      "My Website - Analytics & Stats!",
			wantError: nil,
		},
		{
			name:      "empty string",
			input:     "",
			want:      "",
			wantError: ErrSiteNameTooLong,
		},
		{
			name:      "only whitespace",
			input:     "   ",
			want:      "",
			wantError: ErrSiteNameTooLong,
		},
		{
			name:      "too long - 101 characters",
			input:     "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			want:      "",
			wantError: ErrSiteNameTooLong,
		},
		{
			name:      "unicode characters",
			input:     "My Website 🚀",
			want:      "My Website 🚀",
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSiteName(tt.input)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("ValidateSiteName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateSiteName() = %v, want %v", got, tt.want)
			}
		})
	}
}
