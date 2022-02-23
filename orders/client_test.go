package orders

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

func TestClient_RegisterOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test order validation", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
		}

		_, _, err := RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ReturnURL: cannot be blank"))
	})

	t.Run("Validate return URL", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "wrong\\localhost:6379",
		}

		_, _, err := RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must be a valid URL"))
	})

	t.Run("Validate order number", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
		}

		_, _, err := RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: the length must be between 1 and 30"))
	})

	t.Run("Validate failUrl", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
			FailURL:     "wrong\\localhost:6379",
		}

		_, _, err := RegisterOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("FailURL: must be a valid URL"))
	})

	t.Run("Test Register Order response Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
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

		response, _, err := RegisterOrder(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"OrderId": ContainSubstring("70906e55"),
			"FormUrl": ContainSubstring("application_context"),
		})))

	})
}

func TestClient_RegisterPreAuthOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test order validation", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
		}

		_, _, err := RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ReturnURL: cannot be blank"))
	})

	t.Run("Validate return URL", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "123",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "wrong\\localhost:6379",
		}

		_, _, err := RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must be a valid URL"))
	})

	t.Run("Validate order number", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
		}

		_, _, err := RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: the length must be between 1 and 30"))
	})

	t.Run("Validate failUrl", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://api-sberbank",
			FailURL:     "wrong\\localhost:6379",
		}

		_, _, err := RegisterOrderPreAuth(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("FailURL: must be a valid URL"))
	})

	t.Run("Test Preauth Register Order response Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.RegisterPreAuth, func(w http.ResponseWriter, r *http.Request) {
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

		response, _, err := RegisterOrderPreAuth(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"OrderId": ContainSubstring("70906e55"),
			"FormUrl": ContainSubstring("application_context"),
		})))
	})
}

func TestClient_register(t *testing.T) {
	RegisterTestingT(t)
	//t.Run("Trigger register error on NewRequest", func(t *testing.T) {
	//	newServer := server.NewServer()
	//	defer newServer.Teardown()
	//	prepareClient(newServer.URL)
	//
	//	order := Order{
	//		OrderNumber: "1234567890123456",
	//		Amount:      100,
	//		Description: "Test",
	//		ReturnURL:   "https://localhost",
	//	}
	//
	//	_, _, err := register(context.Background(), "wrong\\localhost:6379", order)
	//	// We dont care what underlying error happened
	//	Expect(err).To(HaveOccurred())
	//})

	t.Run("Trigger register error on Do", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
			Description: "Test",
			ReturnURL:   "https://localhost",
		}
		client := Client{
			API: acquiring.GetAPI(),
		}
		_, _, err := client.register(context.Background(), endpoints.Register, order)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

}

func TestClient_Deposit(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate empty deposit order number", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			Amount: 100,
		}

		_, _, err := Deposit(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))
	})

	t.Run("Validate deposit order number", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      100,
		}

		_, _, err := Deposit(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))
	})

	t.Run("Test Deposit response Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Deposit, func(w http.ResponseWriter, r *http.Request) {
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

		response, _, err := Deposit(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode": Equal(0),
		})))
	})

	t.Run("Test deposit do", func(t *testing.T) {
		t.Run("Trigger register error on Do", func(t *testing.T) {
			newServer := server.NewServer()
			defer newServer.Teardown()
			prepareClient(newServer.URL)

			newServer.Mux.HandleFunc(endpoints.Deposit, func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			})
			order := Order{
				OrderNumber: "1234567890123456",
				Amount:      100,
			}

			_, _, err := Deposit(context.Background(), order)
			// We dont care what underlying error happened
			Expect(err).To(HaveOccurred())
		})
	})

	t.Run("Test deposit with fail on NewRestRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      100,
		}

		_, _, err := Deposit(context.Background(), order)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})
}

func TestClient_ReverseOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate reverse order number", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err := ReverseOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

	})

	t.Run("Test ReverseOrder Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Reverse, func(w http.ResponseWriter, r *http.Request) {
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

		response, _, err := ReverseOrder(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
		})))
	})

	t.Run("Test ReverseOrder 404", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
		}
		_, _, err := ReverseOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("404"))
	})

	t.Run("Test ReverseOrder NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)
		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
		}
		_, _, err := ReverseOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("Test ReverseOrder Do", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Reverse, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		order := Order{
			OrderNumber: "1234567890123456",
		}

		_, _, err := ReverseOrder(context.Background(), order)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})
}

func TestClient_RefundOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate refund order", func(t *testing.T) {

		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "123",
			Amount:      0,
		}

		_, _, err := RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("refund amount should be more"))

		order = Order{
			OrderNumber: "1234567890123456789012345678901",
			Amount:      1,
		}

		_, _, err = RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

		order = Order{
			OrderNumber: "",
			Amount:      1,
		}

		_, _, err = RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))
	})

	t.Run("Test RefundOrder Mapping", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.Refund, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderResponse{
				ErrorCode:    0,
				ErrorMessage: "Успешно",
			})
		})

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
			Amount:      1,
		}

		response, _, err := RefundOrder(context.Background(), order)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
		})))
	})

	t.Run("Test RefundOrder NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
			Amount:      1,
		}
		_, _, err := RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("404"))
	})

	t.Run("Test RefundOrder NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)
		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
			Amount:      1,
		}
		_, _, err := RefundOrder(context.Background(), order)
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
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      1,
		}

		_, _, err := RefundOrder(context.Background(), order)
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
	})
}

func TestClient_ValidateRefundOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("", func(t *testing.T) {
		order := Order{
			OrderNumber: "123",
			Amount:      1,
		}
		err := validateRefundOrder(order)
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_GetOrderStatus(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate order status", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "1234567890123456789012345678901",
		}

		_, _, err := GetOrderStatus(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber is too long"))

		order = Order{
			OrderNumber: "",
		}

		_, _, err = RefundOrder(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orderNumber cant be empty"))
	})

	t.Run("Test GetOrderStatus NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
			Amount:      1,
		}
		_, _, err := GetOrderStatus(context.Background(), order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("404"))
	})

	t.Run("Test GetOrderStatus NewRequest", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		order := Order{
			OrderNumber: "9231a838-ac68-4a3e",
			Amount:      1,
		}
		_, _, err := GetOrderStatus(context.Background(), order)
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("GetOrderStatus Refund Do", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.GetOrderStatusExtended, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      1,
		}

		_, _, err := GetOrderStatus(context.Background(), order)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("GetOrderStatus is working as expected", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()
		prepareClient(newServer.URL)

		newServer.Mux.HandleFunc(endpoints.GetOrderStatusExtended, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderStatusResponse{
				OrderNumber: "70906e55-7114-41d6-8332-4609dc6590f4",
				ErrorCode:   0,
			})
		})
		order := Order{
			OrderNumber: "1234567890123456",
			Amount:      1,
		}

		_, _, err := GetOrderStatus(context.Background(), order)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_ValidateOrder(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test order validator pass", func(t *testing.T) {

		order := Order{
			ReturnURL:   "http://api-sberabcnk.com",
			OrderNumber: "1234",
			FailURL:     "http://123",
		}
		err := order.Validate()
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test order validator fail", func(t *testing.T) {

		order := Order{
			ReturnURL:   "",
			OrderNumber: "1234",
			FailURL:     "http://123",
		}
		err := order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ReturnURL: cannot be blank"))

		order = Order{
			ReturnURL:   "http\\api-sberabcnk.com",
			OrderNumber: "1234",
			FailURL:     "http://google.com",
		}
		err = order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ReturnURL: must be a valid URL."))

		order = Order{
			ReturnURL:   "http:\\api-sberabcnk.com",
			OrderNumber: "1234",
			FailURL:     "http:\\google.com",
		}
		err = order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("FailURL: must be a valid URL; ReturnURL: must be a valid URL."))

		order = Order{
			ReturnURL:   "http://api-sberabcnk.com",
			OrderNumber: "1234",
			FailURL:     "http\\fail-url",
		}
		err = order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("FailURL: must be a valid URL"))

		order = Order{
			ReturnURL:   "http://api-sberabcnk.com",
			OrderNumber: "",
			FailURL:     "http://fail-url",
		}
		err = order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: cannot be blank"))

		order = Order{
			ReturnURL:   "http://api-sberabcnk.com",
			OrderNumber: "123123123123123123213213123123123213213213213213213213213213123123123123123",
			FailURL:     "http://fail-url",
		}
		err = order.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: the length must be between 1 and 30"))
	})
}
