package clientip

import (
	"fmt"
	"net"
	"net/netip"
	"strings"
)

type Resolver struct {
	trusted []netip.Prefix
}

func NewResolver(trustedCIDRs []string) (*Resolver, error) {
	resolver := &Resolver{trusted: make([]netip.Prefix, 0, len(trustedCIDRs))}
	for _, raw := range trustedCIDRs {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		prefix, err := parseTrustedPrefix(value)
		if err != nil {
			return nil, err
		}
		resolver.trusted = append(resolver.trusted, prefix)
	}
	return resolver, nil
}

func MustNewResolver(trustedCIDRs []string) *Resolver {
	resolver, err := NewResolver(trustedCIDRs)
	if err != nil {
		panic(err)
	}
	return resolver
}

func GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {
	return MustNewResolver(nil).GetClientIP(xForwardedFor, xRealIP, remoteAddr)
}

func (r *Resolver) GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {
	remoteHost := hostFromRemoteAddr(remoteAddr)
	remoteIP, remoteOK := parseIP(remoteHost)
	if !remoteOK {
		return remoteHost
	}

	if r == nil || !r.isTrusted(remoteIP) {
		return remoteIP.String()
	}

	if ip, ok := r.clientFromForwardedFor(xForwardedFor); ok {
		return ip.String()
	}
	if ip, ok := parseIP(xRealIP); ok {
		return ip.String()
	}
	return remoteIP.String()
}

func parseTrustedPrefix(value string) (netip.Prefix, error) {
	if strings.Contains(value, "/") {
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			return netip.Prefix{}, fmt.Errorf("invalid trusted proxy CIDR %q: %w", value, err)
		}
		return prefix.Masked(), nil
	}

	addr, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("invalid trusted proxy address %q: %w", value, err)
	}
	if addr.Is4() {
		return netip.PrefixFrom(addr, 32), nil
	}
	return netip.PrefixFrom(addr, 128), nil
}

func (r *Resolver) clientFromForwardedFor(value string) (netip.Addr, bool) {
	parts := strings.Split(value, ",")
	var leftmostValid netip.Addr
	hasLeftmostValid := false

	// Match Nginx real_ip_recursive behavior: from the trusted proxy chain,
	// choose the last non-trusted address as the client.
	for i := len(parts) - 1; i >= 0; i-- {
		ip, ok := parseIP(parts[i])
		if !ok {
			continue
		}
		leftmostValid = ip
		hasLeftmostValid = true
		if !r.isTrusted(ip) {
			return ip, true
		}
	}

	return leftmostValid, hasLeftmostValid
}

func (r *Resolver) isTrusted(ip netip.Addr) bool {
	if r == nil {
		return false
	}
	for _, prefix := range r.trusted {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func hostFromRemoteAddr(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}
	return strings.TrimSpace(remoteAddr)
}

func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return netip.Addr{}, false
	}
	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	}
	addr, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Addr{}, false
	}
	if addr.Is4In6() {
		addr = netip.AddrFrom4(addr.As4())
	}
	return addr, true
}
