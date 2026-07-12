package clientip

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolverIgnoresForwardedHeadersFromUntrustedRemote(t *testing.T) {
	resolver := MustNewResolver([]string{"10.0.0.0/8"})

	got := resolver.GetClientIP("198.51.100.9", "198.51.100.10", "203.0.113.5:12345")

	require.Equal(t, "203.0.113.5", got)
}

func TestResolverUsesForwardedHeaderFromTrustedRemote(t *testing.T) {
	resolver := MustNewResolver([]string{"10.0.0.0/8"})

	got := resolver.GetClientIP("198.51.100.9", "", "10.1.2.3:12345")

	require.Equal(t, "198.51.100.9", got)
}

func TestResolverUsesLastNonTrustedForwardedHop(t *testing.T) {
	resolver := MustNewResolver([]string{"10.0.0.0/8", "192.168.0.0/16"})

	got := resolver.GetClientIP("198.51.100.9, 192.168.1.5, 10.1.2.3", "", "10.9.8.7:12345")

	require.Equal(t, "198.51.100.9", got)
}

func TestResolverFallsBackToXRealIPFromTrustedRemote(t *testing.T) {
	resolver := MustNewResolver([]string{"10.0.0.0/8"})

	got := resolver.GetClientIP("", "198.51.100.11", "10.1.2.3:12345")

	require.Equal(t, "198.51.100.11", got)
}

func TestNewResolverRejectsInvalidCIDR(t *testing.T) {
	_, err := NewResolver([]string{"not-a-cidr"})

	require.Error(t, err)
}
