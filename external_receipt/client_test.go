package receipt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	server "github.com/helios-ag/sberbank-acquiring-go/testing"
	. "github.com/onsi/gomega"
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

func TestClient_GetExternalReceipt(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test External Receipt Validation Request", func(t *testing.T) {
		prepareClient("localhost")

		receipt := ExternalReceiptRequest{}

		_, _, err := GetExternalReceipt(context.Background(), receipt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("userName and Password and mdOrder and Receipt are required"))

	})

	t.Run("GetExternalReceipt error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		receiptRequest := ExternalReceiptRequest{
			UserName: "1234567890123456789018901",
			Password: "test",
			Language: nil,
			MdOrder:  "test",
			Receipt: &Receipt{
				PaymentType: 1,
			},
		}

		_, _, err := GetExternalReceipt(context.Background(), receiptRequest)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("GetExternalReceipt fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub
		receipt := ExternalReceiptRequest{
			UserName: "1234567890123456789018901",
		}

		_, _, err := GetExternalReceipt(context.Background(), receipt)
		// We don't care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("GetExternalReceipt is ok", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderStatusResponse{
				OrderNumber: "70906e55-7114-41d6-8332-4609dc6590f4",
				ErrorCode:   0,
			})
		})

		receipt := ExternalReceiptRequest{
			UserName: "1234567890123456789018901",
			Password: "test",
			MdOrder:  "test",
			Receipt: &Receipt{
				PaymentType: 1,
			},
		}

		_, _, err := GetExternalReceipt(context.Background(), receipt)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}
