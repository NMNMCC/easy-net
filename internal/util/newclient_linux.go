//go:build linux || android
// +build linux android

package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func NewHTTPClient(link string) *http.Client {
	if link == "" {
		return &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	// Create a dialer that binds every socket to the given interface.
	dialer := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			var bindErr error
			err := c.Control(func(fd uintptr) {
				// Use unix.BindToDevice for clarity; equivalent to SetsockoptString(SO_BINDTODEVICE).
				bindErr = unix.BindToDevice(int(fd), link)
			})
			if err != nil {
				return err
			}
			return bindErr
		},
	}

	// Ensure DNS lookups also go out through the bound interface.
	// We do this by providing a custom Resolver whose Dial uses the same bound dialer.
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Network will be udp/udp4/udp6 or tcp/tcp4/tcp6 per Go resolver needs.
			return dialer.DialContext(ctx, network, address)
		},
	}

	// Wrap DialContext to resolve hostnames via the above resolver and then dial the IP
	// using the bound dialer. This guarantees all traffic (DNS + TCP) is bound to the interface.
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		// If address is already an IP, dial directly.
		if ip := net.ParseIP(host); ip != nil {
			return dialer.DialContext(ctx, network, address)
		}

		// Resolve via our resolver (bound to interface)
		ips, err := resolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no A/AAAA records for %s", host)
		}

		// Prefer family according to network hint (tcp4/tcp6)
		prefer4 := strings.Contains(network, "4")
		prefer6 := strings.Contains(network, "6")

		var candidates []net.IP
		if prefer4 {
			for _, ip := range ips {
				if ip.To4() != nil {
					candidates = append(candidates, ip)
				}
			}
		} else if prefer6 {
			for _, ip := range ips {
				if ip.To4() == nil { // treat as v6
					candidates = append(candidates, ip)
				}
			}
		}
		// If no candidates matched preference, try all in order.
		if len(candidates) == 0 {
			candidates = ips
		}

		var lastErr error
		for _, ip := range candidates {
			conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}
		if lastErr == nil {
			lastErr = fmt.Errorf("failed to dial any resolved IP for %s", host)
		}
		return nil, lastErr
	}

	// Clone default transport to inherit sane defaults (TLS, proxies, etc.)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = dialContext

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}
