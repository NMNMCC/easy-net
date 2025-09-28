package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"nmnm.cc/easy-net/internal/log"
)

var loginLogger = log.New("auth/login")

func NewLoginReq(
	base,
	userid, password string,
) (*http.Request, error) {
	url, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	url.Path = "/quickauth.do"

	query := url.Query()

	query.Set("userid", "756"+userid)
	query.Set("passwd", password)

	query.Set("timestamp", fmt.Sprint(time.Now().Unix()))
	query.Set("uuid", uuid.NewString())

	url.RawQuery = query.Encode()

	return http.NewRequest(http.MethodGet, url.String(), io.NopCloser(&bytes.Buffer{}))
}

type LoginRes struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	LogoutGoURL       string `json:"logoutgourl"`
	MacChange         bool   `json:"macChange"`
	GroupID           int    `json:"groupId"`
	PasswdPolicyCheck bool   `json:"passwdPolicyCheck"`
	DropLogCheck      string `json:"dropLogCheck"`
	UserID            string `json:"userId"`
}

func RandomLoginReq(base, password string) (*http.Request, error) {
	userid := RandomUserid()

	return NewLoginReq(base, userid, password)
}

type LoginConfig struct {
	Base     string
	Link     string
	UserID   string
	Password string
}

func Login(cfg *LoginConfig) error {
	client := NewClient(cfg.Link)

	req, err := NewLoginReq(cfg.Base, cfg.UserID, cfg.Password)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	loginLogger.Info("logging in", "url", req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http client failed to do request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", res.StatusCode, cfg.UserID)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var data LoginRes
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}
	if data.Code != "0" {
		loginLogger.Error("login failed", "userid", cfg.UserID, "code", data.Code, "message", data.Message)
		return errors.New(string(lo.Must(json.Marshal(data))))
	}

	loginLogger.Info("login succeeded", "userid", cfg.UserID, "message", data.Message)
	return nil
}
