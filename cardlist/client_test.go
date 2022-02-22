package cardlist

import (
	"context"
	"encoding/json"
	"fmt"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	server "github.com/helios-ag/sberbank-acquiring-go/testing"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

var NewRestRequestStub = func(
	c *acquiring.Client,
	ctx context.Context,
	method,
	urlPath string,
	data map[string]string,
	jsonParams map[string]string) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func prepareClient(URL string) {
	cfg := acquiring.ClientConfig{
		UserName:           "test-api",
		Currency:           currency.RUB,
		Password:           "test",
		Language:           "ru",
		SessionTimeoutSecs: 1200,
		SandboxMode:        true,
	}
	acquiring.SetConfig(cfg)
	acquiring.WithEndpoint(URL)
}

func TestClient_UpdateSSLCardList(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Trigger UpdateSSLCardList error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.UpdateSSLCardList, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		_, _, err := UpdateSSLCardList(context.Background(), "123123", nil)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test UpdateSSLCardList with fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		_, _, err := UpdateSSLCardList(context.Background(), "123123", nil)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("UpdateSSLCardList is working as expected", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.UpdateSSLCardList, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderStatusResponse{
				OrderNumber: "70906e55-7114-41d6-8332-4609dc6590f4",
				ErrorCode:   0,
			})
		})

		_, _, err := UpdateSSLCardList(context.Background(), "123123", nil)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}
