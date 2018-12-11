package acquiring

import (
	"context"
	"encoding/json"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Test Environment
type Server struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	Client *Client
}

func (server *Server) Teardown() {
	server.Server.Close()
	server.Server = nil
	server.Mux = nil
	server.Client = nil
}

func getCfg() *ClientConfig {
	cfg := ClientConfig{
		UserName:           "sb-api",
		Currency:           currency.RUB,
		Password:           "sb",
		Language:           "ru",
		SessionTimeoutSecs: 1200,
		SandboxMode:        false,
	}
	return &cfg
}

func newServer() Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, _ := NewClient(
		getCfg(),
		WithEndpoint(server.URL),
		WithToken("token"),
	)
	return Server{
		Server: server,
		Mux:    mux,
		Client: client,
	}
}

func TestNewClient(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test trailing slashes remove", func(t *testing.T) {
		client, _ := NewClient(getCfg(), WithEndpoint("http://api-sberbank///"))
		if strings.HasSuffix(client.Config.endpoint, "/") {
			t.Fatalf("endpoint has trailing slashes: %q", client.Config.endpoint)
		}
	})
	t.Run("Test getting error response", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				ErrorCode:    4,
				ErrorMessage: "Доступ запрещён.",
			})
		})

		ctx := context.Background()
		_, err := server.Client.NewRestRequest(ctx, "GET", endpoints.Register, nil, nil)
		Expect(err).ShouldNot(HaveOccurred())
	})
}

func TestClientDo(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test client do with external api", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(schema.Response{
				ErrorMessage: "Доступ запрещён.",
				ErrorCode:    5,
			})
		})

		ctx := context.Background()
		request, _ := server.Client.NewRestRequest(ctx, http.MethodGet, endpoints.Register, nil, nil)
		res, err := server.Client.Do(request, nil)
		Expect(err).To(HaveOccurred())
		body, _ := ioutil.ReadAll(res.Body)
		Expect(body).Should(ContainSubstring("\"errorCode\":5"))
	})
}

func TestNewRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Sandbox enabled mode", func(t *testing.T) {
		client, _ := NewClient(
			&ClientConfig{
				UserName:           "sb-api",
				Currency:           currency.RUB,
				Password:           "sb",
				Language:           "ru",
				SessionTimeoutSecs: 1200,
				SandboxMode:        true,
			},
		)

		ctx := context.Background()
		res, err := client.NewRestRequest(ctx, "GET", endpoints.Register, nil, nil)
		Expect(err).ShouldNot(HaveOccurred())
		u, _ := url.Parse(APISandboxURI)
		Expect(res.URL.Host).To(ContainSubstring(u.Host))
	})

	t.Run("Add rest to path", func(t *testing.T) {
		client, _ := NewClient(
			&ClientConfig{
				UserName:           "sb-api",
				Currency:           currency.RUB,
				Password:           "sb",
				Language:           "ru",
				SessionTimeoutSecs: 1200,
				SandboxMode:        true,
			},
		)
		ctx := context.Background()
		_, err := client.NewRequest(ctx, http.MethodGet, endpoints.Register, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("path contains rest request, use NewRestRequest instead"))
	})
	t.Run("Invalid config test", func(t *testing.T) {
		_, err := NewClient(
			&ClientConfig{},
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to validate given config"))
	})

	t.Run("Empty config test", func(t *testing.T) {
		_, err := NewClient(nil)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("passed in config cannot be nil"))
	})
}
