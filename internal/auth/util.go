package auth

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"time"

	ping "github.com/go-ping/ping"

	"nmnm.cc/easy-net/internal/log"
	"nmnm.cc/easy-net/internal/util"
)

var utilLogger = log.New("auth/util")

func TestConnection(link string) (ok bool) {
	// Default target resolver: 9.9.9.9 (Quad9)
	const target = "9.9.9.9"

	utilLogger.Info("testing connection with ICMP/UDP ping", "target", target, "link", link)

	pinger, err := ping.NewPinger(target)
	if err != nil {
		utilLogger.Warn("connection test failed: create pinger", "error", err)
		return false
	}

	// Use unprivileged mode to avoid requiring CAP_NET_RAW; uses UDP fallback.
	// This works for most environments but can be blocked by some firewalls.
	pinger.SetPrivileged(false)
	// Explicitly prefer UDP in unprivileged mode for clarity.
	pinger.SetNetwork("udp")

	// If a specific link (interface) is provided, attempt to set the source IP
	// to ensure packets egress via the bound interface.
	if link != "" {
		ifi, ierr := net.InterfaceByName(link)
		if ierr != nil {
			utilLogger.Warn("failed to find interface for link", "link", link, "error", ierr)
		} else {
			addrs, aerr := ifi.Addrs()
			if aerr != nil {
				utilLogger.Warn("failed to get interface addresses", "link", link, "error", aerr)
			} else {
				var srcIP string
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil {
						continue
					}
					// Prefer IPv4 for 9.9.9.9
					if v4 := ip.To4(); v4 != nil {
						srcIP = v4.String()
						break
					}
				}
				if srcIP != "" {
					pinger.Source = srcIP
					utilLogger.Info("bound ping source IP", "source", srcIP)
				} else {
					utilLogger.Warn("no IPv4 address found on interface to bind source", "link", link)
				}
			}
		}
	}

	// Tuning: small, quick check.
	pinger.Count = 3
	pinger.Timeout = 5 * time.Second
	pinger.Interval = 800 * time.Millisecond

	if err := pinger.Run(); err != nil {
		utilLogger.Warn("connection test failed: run pinger", "error", err)
		return false
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		utilLogger.Warn("connection test failed: no replies", "sent", stats.PacketsSent, "recv", stats.PacketsRecv)
		return false
	}

	utilLogger.Info("connection test succeeded", "sent", stats.PacketsSent, "recv", stats.PacketsRecv, "loss", fmt.Sprintf("%.1f%%", stats.PacketLoss))
	return true
}

var (
	ErrExpectRedirect    = fmt.Errorf("expect redirection")
	ErrExpectRedirectURL = fmt.Errorf("expect redirect URL")
)

func FindPortal(host, link string) (string, error) {
	// Use the bound client so detection traffic goes out through the specified link.
	client := util.NewHTTPClient(link)

	utilLogger.Info("finding portal", "host", host)
	u, _ := url.Parse("http://" + host)
	res, err := client.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("failed to get response from host: %w", err)
	}
	if res.Request.URL.String() == u.String() {
		return "", ErrExpectRedirect
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	re := regexp.MustCompile(`window\.location\.href=\"(.*)"`)

	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", ErrExpectRedirectURL
	}

	p, _ := url.Parse(matches[1])

	final := res.Request.URL.ResolveReference(p).String()

	utilLogger.Info("found portal", "url", final)
	return final, nil
}

// [04] Institute
// [24] Year of Admission
// [04] Class
// [01] Number
func RandomUserid() string {
	institute := rand.Intn(8) + 1
	year := rand.Intn(3) + time.Now().Year()%100 - 3
	class := rand.Intn(10) + 1
	number := rand.Intn(30) + 1

	return fmt.Sprintf("%02d%02d%02d%02d", institute, year, class, number)
}

// NewClient is implemented in platform-specific files.
