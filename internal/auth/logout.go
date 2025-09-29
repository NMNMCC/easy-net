package auth

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"nmnm.cc/easy-net/internal/log"
	"nmnm.cc/easy-net/internal/util"
)

var logoutLogger = log.New("auth/logout")

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

type LogoutConfig struct {
	Base   string
	Link   string
	UserID string
}

func Logout(cfg *LogoutConfig) error {
	client := util.NewHTTPClient(cfg.Link)

	req, err := NewLogoutReq(cfg.Base, cfg.UserID)
	if err != nil {
		return fmt.Errorf("failed to create new logout request: %w", err)
	}

	logoutLogger.Info("logging out", "userid", cfg.UserID)
	if _, err := client.Do(req); err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	logoutLogger.Info("logged out", "userid", cfg.UserID)
	return nil
}
