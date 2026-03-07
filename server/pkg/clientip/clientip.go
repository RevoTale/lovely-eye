package clientip

import (
	"net"
	"strings"
)

func GetClientIP(xForwardedFor, xRealIP, remoteAddr string) string {
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	if xRealIP != "" {
		return xRealIP
	}

	if remoteAddr != "" {
		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			return host
		}
		return remoteAddr
	}

	return ""
}
