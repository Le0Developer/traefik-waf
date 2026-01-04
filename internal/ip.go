package internal

import (
	"net/http"
	"strings"
)

func (i *Instance) getRemoteIP(r *http.Request) string {
	if i.cfg.XffCount >= 0 {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			parts := strings.SplitN(xff, ",", 2)
			return strings.TrimSpace(parts[len(parts)-i.cfg.XffCount])
		}
	}

	ip := r.RemoteAddr
	if colonPos := strings.LastIndex(ip, ":"); colonPos != -1 {
		ip = ip[:colonPos]
	}
	return ip
}
