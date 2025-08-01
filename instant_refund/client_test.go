package instant_refund

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
	. "github.com/onsi/gomega/gstruct"
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

func TestClient_InstantRefund(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate instant refund", func(t *testing.T) {

		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		orderNumber := "9231a838-ac68-4a3e"
		refundRequest := InstantRefundRequest{
			OrderNumber: &orderNumber,
			UserName:    "user",
			Password:    "password",
			Amount:      0,
		}

		_, _, err := InstantRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Amount: cannot be blank."))
		orderNumber = "9231a838-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e"
		refundRequest = InstantRefundRequest{
			UserName:    "user",
			Password:    "password",
			OrderNumber: &orderNumber,
			Amount:      1,
		}

		_, _, err = InstantRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: the length must be between 1 and 30."))

		refundRequest = InstantRefundRequest{
			UserName:    "user",
			Password:    "password",
			OrderNumber: nil,
			Amount:      1,
		}

		_, _, err = InstantRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: cannot be blank."))
	})

	t.Run("Test InstantRefund Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.InstantRefund, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.InstantRefundResponse{
				ErrorCode:    0,
				ErrorMessage: "Успешно",
			})
		})
		orderNumber := "9231a838-ac68-4a3e"
		request := InstantRefundRequest{
			UserName:    "user",
			Password:    "password",
			OrderNumber: &orderNumber,
			Amount:      1,
		}

		response, _, err := InstantRefund(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
		})))
	})

	t.Run("Test InstantRefund NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		orderNumber := "9231a838-ac68-4a3e"
		request := InstantRefundRequest{
			OrderNumber: &orderNumber,
			UserName:    "user",
			Password:    "password",
			Amount:      1,
		}
		_, _, err := InstantRefund(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("404"))
	})

	t.Run("Test InstantRefund NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)
		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		orderNumber := "9231a838-ac68-4a3e"

		request := InstantRefundRequest{
			OrderNumber: &orderNumber,
			Amount:      1,
		}
		_, _, err := InstantRefund(context.Background(), request)
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("Test Refund Do", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Refund, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		orderNumber := "9231a838-ac68-4a3e"

		request := InstantRefundRequest{
			OrderNumber: &orderNumber,
			Amount:      1,
		}

		_, _, err := InstantRefund(context.Background(), request)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
	})
}
