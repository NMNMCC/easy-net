package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

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

func Login(base, userid, password string) error {
	logger := slog.With("component", "login")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := NewLoginReq(base, userid, password)
	if err != nil {
		return err
	}

	logger.Info("logging in", "url", req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		logger.Error("unknown error", "error", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		logger.Error("unknown status code", "status", res.StatusCode)
		return fmt.Errorf("login failed: %s", userid)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("failed to read response body", "error", err)
		return err
	}

	var data LoginRes
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("failed to parse response body", "error", err, "body", string(body))
		return err
	}
	if data.Code != "0" {
		logger.Error("login failed", "userid", userid, "code", data.Code, "message", data.Message)
		return errors.New(string(lo.Must(json.Marshal(data))))
	}

	logger.Info("login succeeded", "userid", userid, "message", data.Message)
	return nil
}
