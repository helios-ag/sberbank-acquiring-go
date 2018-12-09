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

const (
	APIURI        string = "https://securepayments.sberbank.ru"
	APISandboxURI string = "https://3dsec.sberbank.ru"
)

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

type Client struct {
	Config     *ClientConfig
	httpClient *http.Client
}

type Body struct {
	UserName           *string `json:"userName"`
	Token              *string `json:"token"`
	Password           *string `json:"password"`
	Language           *string `json:"language"`
	Currency           *string `json:"currency"`
	SessionTimeoutSecs *int    `json:"sessionTimeoutSecs"`
	JsonParams         []byte  `json:"jsonParams"`
}

type ClientOption func(*Client)

func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.Config.token = token
		client.Config.Password = ""
		client.Config.UserName = ""
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.Config.endpoint = strings.TrimRight(endpoint, "/")
	}
}

func (c *Client) NewRestRequest(ctx context.Context, method, urlPath string, data map[string]string, jsonParams map[string]string) (*http.Request, error) {
	uri := APIURI + urlPath

	if c.Config.SandboxMode {
		uri = APISandboxURI + urlPath
	}

	if c.Config.endpoint != "" {
		uri = c.Config.endpoint + urlPath
	}

	jsonParamsEncoded, err := json.Marshal(jsonParams)

	if err != nil {
		return nil, err
	}

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

	method = "POST"
	reqBodyData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, bytes.NewReader(reqBodyData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) Do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
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
