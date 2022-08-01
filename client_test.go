package sberbank_acquiring_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	server "github.com/helios-ag/sberbank-acquiring-go/testing"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test trailing slashes remove", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "test-api",
			Currency:           currency.RUB,
			Password:           "test",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}
		WithEndpoint("http://api-sberbank///")
		if strings.HasSuffix(cfg.endpoint, "/") {
			t.Fatalf("endpoint has trailing slashes: %q", cfg.endpoint)
		}
	})
	t.Run("Test getting error response", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		testServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.OrderResponse{
				ErrorCode:    4,
				ErrorMessage: "Доступ запрещён.",
			})
		})

		ctx := context.Background()
		_, err := GetAPI().NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		Expect(err).ShouldNot(HaveOccurred())
	})
}

func TestClientDo(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test client do with external api", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		WithEndpoint(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.Response{
				ErrorMessage: "Доступ запрещён.",
				ErrorCode:    5,
			})
		})
		// Override internal reader func with
		reader = func(r io.Reader) (bytes []byte, e error) {
			return nil, errors.New("buf overflow")
		}
		ctx := context.Background()
		request, _ := GetAPI().NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		_, err := GetAPI().Do(request, nil)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("buf overflow"))
		// restore reader
		reader = func(r io.Reader) ([]byte, error) {
			return ioutil.ReadAll(r)
		}
	})

	t.Run("Test response body decode", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		WithEndpoint(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.Response{
				ErrorMessage: "Доступ запрещён.",
				ErrorCode:    5,
			})
		})

		ctx := context.Background()
		request, _ := GetAPI().NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		_, err := GetAPI().Do(request, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test Bad Response Code", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		WithEndpoint(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		ctx := context.Background()
		request, _ := GetAPI().NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		_, err := GetAPI().Do(request, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("sberbank server responded with status code"))
	})
}

func TestErrorFromResponse(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Expect application/json", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("Hello World")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("abc")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("{\"test\": \"test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json header", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application_json")
		body := []byte("{\"test\": test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect Error", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte(`{"errorCode": "5", "errorMessage": "Ошибка"}`)
		err := errorFromResponse(&resp, body)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Ошибка"))
	})

	t.Run("Dont expect Error", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte(`{"errorCode": 0, "errorMessage": ""}`)
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestNewRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Sandbox enabled mode", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "test-api",
			Currency:           currency.RUB,
			Password:           "test",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}
		SetConfig(cfg)

		ctx := context.Background()
		res, err := GetAPI().NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		Expect(err).ShouldNot(HaveOccurred())
		u, _ := url.Parse(APISandboxURI)
		Expect(res.URL.Host).To(ContainSubstring(u.Host))
	})

	t.Run("Add rest to path", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "test-api",
			Currency:           currency.RUB,
			Password:           "test",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}

		SetConfig(cfg)
		ctx := context.Background()
		_, err := GetAPI().NewRequest(ctx, http.MethodGet, endpoints.Register, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("path contains rest request, use NewRestRequest instead"))
	})

	t.Run("Test headers are proper", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "test-api",
			Currency:           currency.RUB,
			Password:           "test",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
			endpoint:           "https://google.com",
		}

		SetConfig(cfg)
		ctx := context.Background()
		req, err := GetAPI().NewRequest(ctx, http.MethodGet, endpoints.GooglePay, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(req.Header.Get("Cache-Control")).To(ContainSubstring("no-cache"))
		Expect(req.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
	})

	t.Run("Trigger NewRequest errors", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "sb-api",
			Currency:           currency.RUB,
			Password:           "sb",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
			endpoint:           APISandboxURI,
		}
		SetConfig(cfg)
		ctx := context.Background()
		// Cyrillic M
		_, err := GetAPI().NewRequest(ctx, "М", endpoints.SamsungPay, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid method"))

		_, err = GetAPI().NewRequest(ctx, "GET", "htt\\wrongUrl", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("Trigger NewRestRequest errors", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "sb-api",
			Currency:           currency.RUB,
			Password:           "sb",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}
		SetConfig(cfg)
		ctx := context.Background()
		// Cyrillic M
		_, err := GetAPI().NewRestRequest(ctx, "М", endpoints.SamsungPay, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid method"))

		_, err = GetAPI().NewRestRequest(ctx, http.MethodGet, "htt\\wrongUrl", nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("NewRestRequest with custom endpoint", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "sb-api",
			Currency:           currency.RUB,
			Password:           "sb",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
			endpoint:           "https://google.com",
		}
		SetConfig(cfg)
		ctx := context.Background()
		_, err := GetAPI().NewRestRequest(ctx, http.MethodGet, "https://google.com", nil, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("NewRestRequest set proper headers", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "sb-api",
			Currency:           currency.RUB,
			Password:           "sb",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}
		SetConfig(cfg)
		ctx := context.Background()
		req, err := GetAPI().NewRestRequest(ctx, http.MethodGet, "https://google.com", nil, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(req.Header.Get("Cache-Control")).To(ContainSubstring("no-cache"))
		Expect(req.Header.Get("Content-Type")).To(ContainSubstring("application/x-www-form-urlencoded"))
	})

}

func TestClientConfigValidation(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test client with empty token and user/password", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "",
			Currency:           currency.RUB,
			token:              "",
			Password:           "",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
		}
		err := cfg.validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Login/Password or Token can't be empty"))
	})

	t.Run("Test client with empty password", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "",
			Currency:           currency.RUB,
			Password:           "",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        false,
			endpoint:           APISandboxURI,
		}
		err := cfg.validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Login/Password or Token can't be empty"))
	})

	t.Run("Test client with wrong endpoint", func(t *testing.T) {
		cfg := ClientConfig{
			UserName:           "abc",
			Currency:           currency.RUB,
			Password:           "asd",
			Language:           "ru",
			SessionTimeoutSecs: 1200,
			SandboxMode:        true,
			endpoint:           "http//google.com",
		}
		err := cfg.validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse URL"))
	})
}
