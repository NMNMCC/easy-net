package internal

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/samber/lo"
)

func TestConnection() (ok bool) {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Returning http.ErrUseLastResponse prevents the client from
			// following the redirect and returns the last received response.
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get("http://captive.apple.com/hotspot-detect.html")
	if err != nil {
		return false
	}

	if res.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func FindPortal(host string) (string, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	u, _ := url.Parse("http://" + host)
	res, err := client.Get(u.String())
	if err != nil {
		return "", err
	}
	if res.Request.URL.String() == u.String() {
		return "", fmt.Errorf("expect redirection")
	}

	body := string(lo.Must(io.ReadAll(res.Body)))

	re := regexp.MustCompile(`window\.location\.href=\"(.*)"`)

	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find portal URL")
	}

	p, _ := url.Parse(matches[1])

	return res.Request.URL.ResolveReference(p).String(), nil
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
