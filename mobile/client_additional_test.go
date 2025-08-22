package mobile

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	server "github.com/helios-ag/sberbank-acquiring-go/testing"
	. "github.com/onsi/gomega"
)

// prepareTestClient is a helper function to set up the test client
func prepareTestClient(URL string) {
	cfg := acquiring.ClientConfig{
		UserName: "test-api",
		Password: "test",
	}
	acquiring.SetConfig(cfg)
	acquiring.WithEndpoint(URL)
}

func TestAdditional_PayWithSamsungPay(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Test validation error", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareTestClient(testServer.URL)

		_, _, err := PayWithSamsungPay(context.Background(), SamsungPaymentRequest{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))
	})

	t.Run("Test successful request", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareTestClient(testServer.URL)

		expectedResponse := schema.SamsungPaymentResponse{
			Success: true,
		}

		testServer.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedResponse)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test-order",
			Merchant:     "test-merchant",
			PaymentToken: "test-token-123",
		}

		resp, httpResp, err := PayWithSamsungPay(context.Background(), req)

		Expect(err).NotTo(HaveOccurred())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Success).To(BeTrue())
	})

	t.Run("Test API error response", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareTestClient(testServer.URL)

		errorResponse := struct {
			ErrorCode    int    `json:"errorCode"`
			ErrorMessage string `json:"errorMessage"`
		}{
			ErrorCode:    5,
			ErrorMessage: "Invalid parameters",
		}

		testServer.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test-order",
			Merchant:     "test-merchant",
			PaymentToken: "test-token-123",
		}

		_, httpResp, err := PayWithSamsungPay(context.Background(), req)

		Expect(err).To(HaveOccurred())
		Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
	})
}
