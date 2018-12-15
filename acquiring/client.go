package acquiring

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// URLS for API endpoints
const (
	APIURI        string = "https://securepayments.sberbank.ru"
	APISandboxURI string = "https://3dsec.sberbank.ru"
)

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

// Client is a client to SB API
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
func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.Config.token = token
		client.Config.Password = ""
		client.Config.UserName = ""
	}
}

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.Config.endpoint = strings.TrimRight(endpoint, "/")
	}
}

// NewRestRequest creates an HTTP request against the API with 'rest' in path. The returned request
// is assigned with ctx and has all necessary headers set (auth, user agent, etc.).
func (c *Client) NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string, jsonParams map[string]string) (*http.Request, error) {
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
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, data interface{}) (*http.Request, error) {
	if strings.Contains(urlPath, "rest") {
		return nil, fmt.Errorf("path contains rest request, use NewRestRequest instead")
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
	return ioutil.ReadAll(r)
}

// Do performs an HTTP request against the API.
func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
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
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

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

	if _, err := url.Parse(c.endpoint); err != nil {
		return fmt.Errorf("unable to parse URL: %v", err)
	}

	return nil
}

// NewClient creates a new client.
func NewClient(cfg *ClientConfig, options ...ClientOption) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("passed in config cannot be nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("unable to validate given config: %v", err)
	}

	client := &Client{
		Config:     cfg,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}
