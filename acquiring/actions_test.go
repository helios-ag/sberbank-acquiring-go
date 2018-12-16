package acquiring

import (
	"context"
	"encoding/json"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"testing"
)

func prepareClient() (*Client, error) {
	cfg := ClientConfig{
		UserName:           "test-api",
		Currency:           currency.RUB,
		Password:           "test",
		Language:           "ru",
		SessionTimeoutSecs: 1200,
		SandboxMode:        true,
	}

	client, err := NewClient(&cfg, WithEndpoint("http://api-sberbank///"))

	return client, err
}

func TestRegisterOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test order validation", func(t *testing.T) {

		client, _ := prepareClient()

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
		}

		_, _, err := client.RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("return url must be set"))
	})

	t.Run("Validate return URL", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "wrong\\localhost:6379",
		}

		_, _, err := client.RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse"))
	})

	t.Run("Validate order number", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
		}

		_, _, err := client.RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("Validate failUrl", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
			FailURL:     "wrong\\localhost:6379",
		}

		_, _, err := client.RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse FailUrl"))
	})

	t.Run("Test Register Order response Mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				OrderId: "70906e55-7114-41d6-8332-4609dc6590f4",
				FormUrl: "https://server/application_context/merchants/test/payment_ru.html?mdOrder=70906e55-7114-41d6-8332-4609dc6590f4",
			})
		})

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://localhost",
		}

		response, _, err := server.Client.RegisterOrder(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"OrderId": ContainSubstring("70906e55"),
			"FormUrl": ContainSubstring("application_context"),
		})))

	})
}

func TestRegisterPreAuthOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test order validation", func(t *testing.T) {

		client, _ := prepareClient()

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
		}

		_, _, err := client.RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("return url must be set"))
	})

	t.Run("Validate return URL", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "wrong\\localhost:6379",
		}

		_, _, err := client.RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse"))
	})

	t.Run("Validate order number", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
		}

		_, _, err := client.RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("Validate failUrl", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
			FailURL:     "wrong\\localhost:6379",
		}

		_, _, err := client.RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse FailUrl"))
	})
	// TODO: Replace with ghttp from gomega
	t.Run("Test Preauth Register Order response Mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.RegisterPreAuth, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				OrderId: "70906e55-7114-41d6-8332-4609dc6590f4",
				FormUrl: "https://server/application_context/merchants/test/payment_ru.html?mdOrder=70906e55-7114-41d6-8332-4609dc6590f4",
			})
		})

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://localhost",
		}

		response, _, err := server.Client.RegisterOrderPreAuth(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"OrderId": ContainSubstring("70906e55"),
			"FormUrl": ContainSubstring("application_context"),
		})))
	})
}

func TestRegister(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Trigger register error on NewRequest", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://localhost",
		}

		_, _, err := server.Client.register(context.Background(), "wrong\\localhost:6379", order)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Trigger register error on Do", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		})
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://localhost",
		}

		_, _, err := server.Client.register(context.Background(), endpoints.Register, order)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

}

func TestDeposit(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate empty deposit order number", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			Amount: 100,
		}

		_, _, err := client.Deposit(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))
	})

	t.Run("Validate deposit order number", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
		}

		_, _, err := client.Deposit(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("Test Deposit response Mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Deposit, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				ErrorCode: 0,
			})
		})

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
		}

		response, _, err := server.Client.Deposit(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode": Equal(0),
		})))
	})

	t.Run("Test deposit do", func(t *testing.T) {
		t.Run("Trigger register error on Do", func(t *testing.T) {
			server := newServer()
			defer server.Teardown()

			server.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			})
			order := Order{
				OrderNumber: "1234567890123456",
				Amount:      100,
			}

			_, _, err := server.Client.Deposit(context.Background(), order)
			// We dont care what underlying error happened
			Expect(err).To(HaveOccurred())
		})
	})
}

func TestReverseOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate reverse order number", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err := client.ReverseOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

	})

	t.Run("Test Reverse order Mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.Deposit, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				ErrorCode:    0,
				ErrorMessage: "Успешно",
			})
		})

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
		}

		response, _, err := server.Client.Deposit(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
		})))
	})
}

func TestRefundOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate refund order", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "123",
			Amount:      0,
		}

		_, _, err := client.RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("refund amount should be more"))

		order = Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      1,
		}

		_, _, err = client.RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

		order = Order{
			OrderNumber: "",
			Amount:      1,
		}

		_, _, err = client.RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))
	})
}

func TestGetOrderStatus(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate order status", func(t *testing.T) {
		client, _ := prepareClient()

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err := client.RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

		order = Order{
			OrderNumber: "",
		}

		_, _, err = client.RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))

	})
}

func TestEnrollment(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate Enrollment PAN", func(t *testing.T) {
		client, _ := prepareClient()

		enrollment := "4111111111111111111111111"

		_, _, err := client.VerifyEnrollment(context.Background(), enrollment)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("pan length shouldn't be less 13 or more 19 symbols"))
	})

	t.Run("Test Enrollment response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.VerifyEnrollment, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.EnrollmentResponse{
				ErrorCode:    0,
				ErrorMessage: "Успешно",
				EmitterName:  "TEST",
				Enrolled:     'Y',
			})
		})

		enrollment := "4111111111111111"

		response, _, err := server.Client.VerifyEnrollment(context.Background(), enrollment)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
			"Enrolled":     Equal(byte('Y')),
		})))
	})
}

func TestBindingCard(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Bind Validate", func(t *testing.T) {
		client, _ := prepareClient()

		binding := Binding{
			bindingID: "",
		}

		_, _, err := client.BindCard(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("bindingId can't be empty"))
	})

	t.Run("Test UnBind Validate", func(t *testing.T) {
		client, _ := prepareClient()

		binding := Binding{
			bindingID: "",
		}

		_, _, err := client.UnBindCard(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("bindingId can't be empty"))
	})

	t.Run("Test validate Extend Binding", func(t *testing.T) {
		client, _ := prepareClient()

		binding := Binding{
			bindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
		}

		_, _, err := client.ExtendBinding(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("new expiry date should have 6 digits"))
	})

	t.Run("Test validate get bindings", func(t *testing.T) {
		client, _ := prepareClient()

		longString := make([]byte, 300)
		_, _, err := client.GetBindings(context.Background(), string(longString), nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("clientId is too long"))
	})

	t.Run("Test Binding response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.BindCard, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    2,
				ErrorMessage: "Binding is active",
			})
		})

		binding := Binding{
			bindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
		}

		response, _, err := server.Client.BindCard(context.Background(), binding)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(2),
			"ErrorMessage": Equal("Binding is active"),
		})))
	})

}

func TestReceiptStatus(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Receipt Status Validation Request", func(t *testing.T) {
		client, _ := prepareClient()

		receipt := ReceiptStatusRequest{}

		_, _, err := client.GetReceiptStatus(context.Background(), receipt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("pass orderNumber"))

		receipt = ReceiptStatusRequest{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err = client.GetReceiptStatus(context.Background(), receipt)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("Test Binding response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.BindCard, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    2,
				ErrorMessage: "Binding is active",
			})
		})

		binding := Binding{
			bindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
		}

		response, _, err := server.Client.BindCard(context.Background(), binding)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(2),
			"ErrorMessage": Equal("Binding is active"),
		})))
	})
}

func TestApplePaymentRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {
		client, _ := prepareClient()

		req := ApplePaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := client.PayWithApplePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = ApplePaymentRequest{
			OrderNumber: "123",
			Merchant:    "123",
		}

		_, _, err = client.PayWithApplePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Apple Payment response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		resp := schema.ApplePaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		server.Mux.HandleFunc(endpoints.ApplePay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := ApplePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := server.Client.PayWithApplePay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

}

func TestGooglePaymentRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {
		client, _ := prepareClient()

		req := GooglePaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := client.PayWithGooglePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		_, _, err = client.PayWithGooglePay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Google Payment response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		resp := schema.GooglePaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		server.Mux.HandleFunc(endpoints.GooglePay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := GooglePaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := server.Client.PayWithGooglePay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})

}

func TestSamsungPaymentRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test apple payment request validation", func(t *testing.T) {
		client, _ := prepareClient()

		req := SamsungPaymentRequest{
			OrderNumber: "123",
		}

		_, _, err := client.PayWithSamsungPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber, merchant and PaymentToken are required"))

		req = SamsungPaymentRequest{
			OrderNumber:  "123",
			Merchant:     "123",
			PaymentToken: "test",
		}

		_, _, err = client.PayWithSamsungPay(context.Background(), req)
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test Samsung Payment response mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		resp := schema.SamsungPaymentResponse{
			Success: true,
		}

		resp.Data.OrderID = "b926351f-a634-49cf-9484-ccb0a3b8cfad"

		server.Mux.HandleFunc(endpoints.SamsungPay, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		})

		req := SamsungPaymentRequest{
			OrderNumber:  "test",
			Merchant:     "test",
			PaymentToken: "test",
		}

		response, _, err := server.Client.PayWithSamsungPay(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Success": Equal(true),
			"Data": MatchFields(IgnoreExtras, Fields{
				"OrderID": Equal("b926351f-a634-49cf-9484-ccb0a3b8cfad"),
			}),
		})))
	})
}
