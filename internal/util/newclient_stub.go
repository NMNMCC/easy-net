//go:build !linux && !android
// +build !linux,!android

package util

import (
	"log/slog"
	"net/http"
	"time"
)

func NewHTTPClient(link string) *http.Client {
	if link != "" {
		slog.Warn("SO_BINDTODEVICE is not supported on this platform; falling back to default client", "link", link)
	}

	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
