package receipt

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


var NewRestRequestStub = func(
	c *acquiring.Client,
	ctx context.Context,
	method,
	urlPath string,
	data map[string]string,
	jsonParams map[string]string) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func TestClient_GetReceiptStatus(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Receipt Status Validation Request", func(t *testing.T) {
		prepareClient("localhost")

		receipt := StatusRequest{}

		_, _, err := GetReceiptStatus(context.Background(), receipt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("pass orderNumber"))

		receipt = StatusRequest{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err = GetReceiptStatus(context.Background(), receipt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("GetReceiptStatus error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetReceiptStatus, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		receipt := StatusRequest{
			OrderNumber: "1234567890123456789018901",
		}

		_, _, err := GetReceiptStatus(context.Background(), receipt)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("GetReceiptStatus fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub
		receipt := StatusRequest{
			OrderNumber: "1234567890123456789018901",
		}

		_, _, err := GetReceiptStatus(context.Background(), receipt)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("GetReceiptStatus is ok", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetReceiptStatus, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderStatusResponse{
				OrderNumber: "70906e55-7114-41d6-8332-4609dc6590f4",
				ErrorCode:   0,
			})
		})

		receipt := StatusRequest{
			OrderNumber: "1234567890123456789018901",
		}

		_, _, err := GetReceiptStatus(context.Background(), receipt)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}
