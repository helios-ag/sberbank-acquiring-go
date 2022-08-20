package sberbank_acquiring_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// URLS for API endpoints
const (
	APIURI        string = "https://securepayments.sberbank.ru"
	APISandboxURI string = "https://3dsec.sberbank.ru"
)

// APIs are the currently supported endpoints.
type APIs struct {
	Api API
	mu  sync.RWMutex
}

var apis APIs

// ClientConfig is used to set client configuration
type ClientConfig struct {
	UserName           string
	Password           string
	Currency           int
	Language           string
	SessionTimeoutSecs int
	endpoint           string
	token              string
	SandboxMode        bool
}

type API interface {
	NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string, jsonParams map[string]string) (*http.Request, error)
	NewRequest(ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error)
	Do(r *http.Request, v interface{}) (*http.Response, error)
}

// Client is a client to SberBank API
type Client struct {
	Config     *ClientConfig
	httpClient *http.Client
}

// Body struct
type Body struct {
	UserName           *string `json:"userName"`
	Token              *string `json:"token"`
	Password           *string `json:"password"`
	Language           *string `json:"language"`
	Currency           *string `json:"currency"`
	SessionTimeoutSecs *int    `json:"sessionTimeoutSecs"`
	JsonParams         []byte  `json:"jsonParams"`
}

// ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithToken configures a Client to use the specified token for authentication.
func WithToken(token string) {
	cfg.token = token
	cfg.Password = ""
	cfg.UserName = ""
}

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) {
	cfg.endpoint = strings.TrimRight(endpoint, "/")
}

// NewRestRequest creates an HTTP request against the API with 'rest' in path. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string, jsonParams map[string]string) (*http.Request, error) {
	return NewRestRequest(c, ctx, method, urlPath, data, jsonParams)
}

var NewRestRequest = func(c *Client, ctx context.Context, method, urlPath string, data map[string]string, jsonParams map[string]string) (*http.Request, error) {
	uri := APIURI + urlPath

	if c.Config.SandboxMode {
		uri = APISandboxURI + urlPath
	}

	if c.Config.endpoint != "" {
		uri = c.Config.endpoint + urlPath
	}

	jsonParamsEncoded, _ := json.Marshal(jsonParams)

	body := url.Values{}
	body.Add("userName", c.Config.UserName)
	body.Add("password", c.Config.Password)
	body.Add("currency", strconv.Itoa(c.Config.Currency))
	body.Add("jsonParams", string(jsonParamsEncoded[:]))
	body.Add("sessionTimeoutSecs", strconv.Itoa(c.Config.SessionTimeoutSecs))

	for key, value := range data {
		body.Add(key, value)
	}
	reqData := body.Encode()
	req, err := http.NewRequest(method, uri, strings.NewReader(reqData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req = req.WithContext(ctx)
	return req, nil
}

// NewRequest creates an HTTP request against the API (mobile payments). The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
// NewRestRequest creates an HTTP request against the API with 'rest' in path. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error) {
	return NewRequest(c, ctx, method, urlPath, data)
}

var NewRequest = func(c *Client, ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error) {
	if strings.Contains(urlPath, "rest") {
		return nil, fmt.Errorf("path contains rest request, use NewRestRequest instead")
	}

	if err := c.Config.validate(); err != nil {
		return nil, err
	}

	uri := APIURI + urlPath

	if c.Config.SandboxMode {
		uri = APISandboxURI + urlPath
	}

	if c.Config.endpoint != "" {
		uri = c.Config.endpoint + urlPath
	}

	reqBodyData, _ := json.Marshal(data)

	req, err := http.NewRequest(method, uri, bytes.NewReader(reqBodyData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)

	return req, nil
}

var reader = func(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// Do perform an HTTP request against the API.
func (c Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := reader(resp.Body)
	if err != nil {
		resp.Body.Close()
		return resp, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		err = errorFromResponse(resp, body)
		if err == nil {
			err = fmt.Errorf("sberbank server responded with status code %d", resp.StatusCode)
		}
		return resp, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return resp, err
}

func errorFromResponse(resp *http.Response, body []byte) error {
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	var respBody schema.Response
	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil
	}
	if respBody.ErrorCode == 0 && respBody.ErrorMessage == "" {
		return nil
	}
	return fmt.Errorf(respBody.ErrorMessage)
}

func (c *ClientConfig) validate() error {
	if (c.Password == "" && c.UserName == "") || c.token != "" {
		return errors.New("Login/Password or Token can't be empty")
	}

	if _, err := url.ParseRequestURI(c.endpoint); err != nil {
		return fmt.Errorf("unable to parse URL: %v", err)
	}

	return nil
}

// newAPI creates a new client.
func newAPI(cfg *ClientConfig, options ...ClientOption) *Client {
	client := &Client{
		Config:     cfg,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func GetAPI(options ...ClientOption) API {
	var api API

	apis.mu.RLock()
	api = apis.Api
	apis.mu.RUnlock()

	if api != nil {
		return api
	}

	return newAPI(&cfg, options...)
}

var cfg ClientConfig

func SetConfig(config ClientConfig) {
	cfg = config
}
