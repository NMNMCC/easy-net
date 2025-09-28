//go:build linux || android
// +build linux android

package auth

import (
	"net"
	"net/http"
	"syscall"
	"time"
)

func NewClient(link string) *http.Client {
	if link == "" {
		return &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	dialer := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			var bindErr error
			err := c.Control(func(fd uintptr) {
				bindErr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, link)
			})
			if err != nil {
				return err
			}
			return bindErr
		},
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
	}

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}
