package internal

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func NewLogoutReq(base, userid string) (*http.Request, error) {
	url, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	url.Path = "quickauthdisconn.do"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, values := range url.Query() {
		if key == "userid" {
			writer.WriteField("userid", userid+"@zk")
			continue
		}
		for _, value := range values {
			writer.WriteField(key, value)
		}
	}

	return http.NewRequest("POST", url.String(), io.NopCloser(body))
}

func Logout(base, userid string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := NewLogoutReq(base, userid)
	if err != nil {
		return err
	}

	if _, err := client.Do(req); err != nil {
		return err
	}

	return nil
}
