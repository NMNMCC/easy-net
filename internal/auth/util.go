package auth

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"syscall"
	"time"
)

func TestConnection() (ok bool) {
	logger := slog.With("component", "test-connection")

	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Returning http.ErrUseLastResponse prevents the client from
			// following the redirect and returns the last received response.
			return http.ErrUseLastResponse
		},
	}

	logger.Info("testing connection with http://captive.apple.com/hotspot-detect.html")
	res, err := client.Get("http://captive.apple.com/hotspot-detect.html")
	if err != nil {
		logger.Warn("connection test failed", "error", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		logger.Warn("connection test failed", "status", res.StatusCode)
		return false
	}

	logger.Info("connection test succeeded")
	return true
}

var (
	ErrExpectRedirect    = fmt.Errorf("expect redirection")
	ErrExpectRedirectURL = fmt.Errorf("expect redirect URL")
)

func FindPortal(host string) (string, error) {
	logger := slog.With("component", "find-portal")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	logger.Info("finding portal", "host", host)
	u, _ := url.Parse("http://" + host)
	res, err := client.Get(u.String())
	if err != nil {
		logger.Error("unknown error", "error", err)
		return "", err
	}
	if res.Request.URL.String() == u.String() {
		logger.Warn("expect redirect", "url", u.String())
		return "", ErrExpectRedirect
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Warn("failed to read response body", "error", err)
		return "", err
	}

	re := regexp.MustCompile(`window\.location\.href=\"(.*)"`)

	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		logger.Warn("expect redirect URL", "body", string(body))
		return "", ErrExpectRedirectURL
	}

	p, _ := url.Parse(matches[1])

	final := res.Request.URL.ResolveReference(p).String()

	logger.Info("found portal", "url", final)
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

func NewClient(link string) *http.Client {
	var client *http.Client
	if link != "" {
		dialer := &net.Dialer{
			Control: func(network, address string, c syscall.RawConn) error {
				var be error
				err := c.Control(func(fd uintptr) {
					be = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, link)
				})

				if err != nil {
					return err
				}
				return be
			},
		}
		transport := &http.Transport{
			DialContext: dialer.DialContext,
		}
		client = &http.Client{
			Timeout:   5 * time.Second,
			Transport: transport,
		}
	} else {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return client
}
