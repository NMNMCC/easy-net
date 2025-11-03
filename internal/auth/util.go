package auth

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

func randomUint32N(min, max uint64) uint64 {
	return rand.Uint64N(max-min+1) + min
}

func RandomUserid(
	// [04] Institute
	instituteMin, instituteMax,
	// [24] Year of Admission
	yearMin, yearMax,
	// [04] Class
	classMin, classMax,
	// [01] Number
	idMin, idMax uint64,
) string {
	institute := randomUint32N(instituteMin, instituteMax)
	year := randomUint32N(yearMin, yearMax) % 100
	class := randomUint32N(classMin, classMax)
	id := randomUint32N(idMin, idMax)

	return fmt.Sprintf("%02d%02d%02d%02d", institute, year, class, id)
}

// 00000000-XXXXXXXX
func ParseRange(r string) (
	// [04] Institute
	instituteMin, instituteMax,
	// [24] Year of Admission
	yearMin, yearMax,
	// [04] Class
	classMin, classMax,
	// [01] Number
	idMin, idMax uint64, err error) {
	parts := strings.SplitN(r, "-", 2)
	min := parts[0]
	max := parts[1]

	iMin, err := strconv.ParseUint(min[0:2], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	iMax, err := strconv.ParseUint(max[0:2], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	yMin, err := strconv.ParseUint(min[2:4], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	yMax, err := strconv.ParseUint(max[2:4], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	cMin, err := strconv.ParseUint(min[4:6], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	cMax, err := strconv.ParseUint(max[4:6], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	dMin, err := strconv.ParseUint(min[6:8], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}
	dMax, err := strconv.ParseUint(max[6:8], 10, 64)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), uint64(0), err
	}

	return iMin, iMax, yMin + 1000, yMax + 1000, cMin, cMax, dMin, dMax, nil
}
