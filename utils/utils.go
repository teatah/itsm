package utils

import (
	"net/http"
	"strings"
)

func GetPort(r *http.Request) string {
	hostParts := strings.Split(r.Host, ":")
	var port string
	if len(hostParts) > 1 {
		port = hostParts[1]
	} else {
		port = "8080"
	}
	return port
}
