package auth

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"nmnm.cc/easy-net/internal/log"
)

var utilLogger = log.New("auth/util")

func TestConnection() (ok bool) {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Returning http.ErrUseLastResponse prevents the client from
			// following the redirect and returns the last received response.
			return http.ErrUseLastResponse
		},
	}

	utilLogger.Info("testing connection with http://captive.apple.com/hotspot-detect.html")
	res, err := client.Get("http://captive.apple.com/hotspot-detect.html")
	if err != nil {
		utilLogger.Warn("connection test failed", "error", err)
		return false
	}
	if res.StatusCode != http.StatusOK {
		utilLogger.Warn("connection test failed", "status", res.StatusCode)
		return false
	}

	utilLogger.Info("connection test succeeded")
	return true
}

var (
	ErrExpectRedirect    = fmt.Errorf("expect redirection")
	ErrExpectRedirectURL = fmt.Errorf("expect redirect URL")
)

func FindPortal(host string) (string, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

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
