package helpers

import (
	"net"
	"net/http"
)

// GetIP returns the real user IP address
func GetIP(r *http.Request) string {
	if xForwardedFor := r.Header.Get("X-FORWARDED-FOR"); len(xForwardedFor) > 0 {
		return xForwardedFor
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
