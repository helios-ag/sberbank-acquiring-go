package mobile

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
	. "github.com/onsi/gomega/gstruct"
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

var NewRequestStub = func(
	c *acquiring.Client,
	ctx context.Context,
	method,
	urlPath string,
	data interface{},
) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func TestClient_PayWithApplePay(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {

		req := ApplePaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := PayWithApplePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = ApplePaymentRequest{
			OrderNumber: "123",
			Merchant:    "123",
		}

		_, _, err = PayWithApplePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Apple Payment response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		resp := schema.ApplePaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		testServer.Mux.HandleFunc(endpoints.ApplePay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := ApplePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := PayWithApplePay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

	t.Run("Test ApplePaymentDo Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		testServer.Mux.HandleFunc(endpoints.ApplePay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := ApplePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithApplePay(context.Background(), req)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test ApplePaymentDo NewRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		oldNewRequest := acquiring.NewRequest
		acquiring.NewRequest = NewRequestStub

		testServer.Mux.HandleFunc(endpoints.ApplePay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := ApplePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithApplePay(context.Background(), req)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
		acquiring.NewRequest = oldNewRequest
	})
}

func TestClient_PayWithGooglePay(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {

		req := GooglePaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := PayWithGooglePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err = PayWithGooglePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Google Payment response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		resp := schema.GooglePaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		testServer.Mux.HandleFunc(endpoints.GooglePay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := PayWithGooglePay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

	t.Run("Test GooglePaymentRequest Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		testServer.Mux.HandleFunc(endpoints.GooglePay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithGooglePay(context.Background(), req)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test GooglePaymentRequest NewRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		oldNewRequest := acquiring.NewRequest
		acquiring.NewRequest = NewRequestStub
		testServer.Mux.HandleFunc(endpoints.GooglePay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithGooglePay(context.Background(), req)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
		acquiring.NewRequest = oldNewRequest
	})
}
func TestClient_PayWithSamsungPay(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := SamsungPaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := PayWithSamsungPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = SamsungPaymentRequest{
			OrderNumber:  "123",
			Merchant:     "123",
			PaymentToken: "test",
		}

		_, _, err = PayWithSamsungPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Samsung Payment response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		resp := schema.SamsungPaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		testServer.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := PayWithSamsungPay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

	t.Run("Test Samsung Payment Request Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		testServer.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithSamsungPay(context.Background(), req)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test GooglePaymentRequest NewRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRequest
		acquiring.NewRequest = NewRequestStub
		testServer.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithSamsungPay(context.Background(), req)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
		acquiring.NewRequest = oldNewRequest
	})
}

func TestClient_PayWithMirPay(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test mir payment request validation", func(t *testing.T) {

		req := MirPayPaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := PayWithMirPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant, IP and PaymentToken are required"))

		req = MirPayPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err = PayWithMirPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test mir Payment response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		resp := schema.MirPayPaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		testServer.Mux.HandleFunc(endpoints.MirPay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := MirPayPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
			IP: "10.10.10.1",
		}

		response, _, err := PayWithMirPay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

	t.Run("Test MirPayment Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()

		testServer.Mux.HandleFunc(endpoints.MirPay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := MirPayPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithMirPay(context.Background(), req)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test MirPayment NewRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		oldNewRequest := acquiring.NewRequest
		acquiring.NewRequest = NewRequestStub
		testServer.Mux.HandleFunc(endpoints.GooglePay, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		req := MirPayPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err := PayWithMirPay(context.Background(), req)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
		acquiring.NewRequest = oldNewRequest
	})
}
