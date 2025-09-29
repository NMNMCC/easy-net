package auth

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"nmnm.cc/easy-net/internal/log"
	"nmnm.cc/easy-net/internal/util"
)

var utilLogger = log.New("auth/util")

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
	institute := rand.Intn(4) + 1
	year := rand.Intn(2) + time.Now().Year()%100 - 3
	class := rand.Intn(10) + 1
	number := rand.Intn(30) + 1

	return fmt.Sprintf("%02d%02d%02d%02d", institute, year, class, number)
}

// NewClient is implemented in platform-specific files.
