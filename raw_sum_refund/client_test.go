package raw_sum_refund

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

func TestClient_RawSumRefund(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate RawSum refund", func(t *testing.T) {

		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		OrderId := "9231a838-ac68-4a3e"
		refundRequest := ProcessRawSumRefundRequest{
			OrderId:  OrderId,
			UserName: "user",
			Password: "password",
			Amount:   0,
		}

		_, _, err := ProcessRawSumRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Amount: cannot be blank."))
		OrderId = "9231a838-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e-ac68-4a3e"
		refundRequest = ProcessRawSumRefundRequest{
			UserName: "user",
			Password: "password",
			OrderId:  OrderId,
			Amount:   1,
		}

		_, _, err = ProcessRawSumRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderId: the length must be between 1 and 30."))

		refundRequest = ProcessRawSumRefundRequest{
			UserName: "user",
			Password: "password",
			OrderId:  "",
			Amount:   1,
		}

		_, _, err = ProcessRawSumRefund(context.Background(), refundRequest)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderId: cannot be blank."))
	})

	t.Run("Test RawSumRefund Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.ProcessRawSumRefund, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.ProcessRawRefundResponse{
				ErrorCode:    0,
				ErrorMessage: "Успешно",
			})
		})
		OrderId := "9231a838-ac68-4a3e"
		request := ProcessRawSumRefundRequest{
			UserName: "user",
			Password: "password",
			OrderId:  OrderId,
			Amount:   1,
		}

		response, _, err := ProcessRawSumRefund(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
		})))
	})

	t.Run("Test RawSumRefund NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		OrderId := "9231a838-ac68-4a3e"
		request := ProcessRawSumRefundRequest{
			OrderId:  OrderId,
			UserName: "user",
			Password: "password",
			Amount:   1,
		}
		_, _, err := ProcessRawSumRefund(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("404"))
	})

	t.Run("Test RawSumRefund NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)
		oldNewRequest := acquiring.NewRestRequest
		defer func() { acquiring.NewRestRequest = oldNewRequest }()

		acquiring.NewRestRequest = NewRestRequestStub

		OrderId := "9231a838-ac68-4a3e"

		request := ProcessRawSumRefundRequest{
			OrderId: OrderId,
			Amount:  1,
		}
		_, _, err := ProcessRawSumRefund(context.Background(), request)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Refund Do", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.ProcessRawSumRefund, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		OrderId := "9231a838-ac68-4a3e"

		request := ProcessRawSumRefundRequest{
			OrderId: OrderId,
			Amount:  1,
		}

		_, _, err := ProcessRawSumRefund(context.Background(), request)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
	})
}
