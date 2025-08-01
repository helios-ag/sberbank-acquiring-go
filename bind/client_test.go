package bind

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

func TestClient_BindCard(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Bind Validate", func(t *testing.T) {
		binding := Binding{
			BindingID: "",
		}

		_, _, err := BindCard(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("BindingID: cannot be blank"))
	})

	t.Run("Test Binding response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.BindCard, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    2,
				ErrorMessage: "Binding is active",
			})
		})

		binding := Binding{
			BindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
		}

		response, _, err := BindCard(context.Background(), binding)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(2),
			"ErrorMessage": Equal("Binding is active"),
		})))
	})

	t.Run("Test UnBind Validate", func(t *testing.T) {
		binding := Binding{
			BindingID: "",
		}

		_, _, err := UnBindCard(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("BindingID: cannot be blank"))
	})

	t.Run("Test validate ExtendBinding with empty value", func(t *testing.T) {
		binding := Binding{
			BindingID: "",
		}

		_, _, err := ExtendBinding(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("BindingID: cannot be blank"))
	})

	t.Run("Test validate ExtendBinding Expiry", func(t *testing.T) {
		binding := Binding{
			BindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
			NewExpiry: "123",
		}
		_, _, err := ExtendBinding(context.Background(), binding)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("NewExpiry: must be in a valid format."))
	})

	t.Run("Test ExtendBinding is ok", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.ExtendBinding, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    2,
				ErrorMessage: "Binding is active",
			})
		})

		binding := Binding{
			BindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
			NewExpiry: "123123",
		}
		_, _, err := ExtendBinding(context.Background(), binding)
		// We don't care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test validate get bindings", func(t *testing.T) {
		longString := make([]byte, 300)

		_, _, err := GetBindings(context.Background(), string(longString), nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("clientId is too long"))
	})

	t.Run("Test Binding response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.BindCard, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    2,
				ErrorMessage: "Binding is active",
			})
		})

		binding := Binding{
			BindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
		}

		response, _, err := BindCard(context.Background(), binding)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(2),
			"ErrorMessage": Equal("Binding is active"),
		})))
	})

}

func TestClient_ValidateBind(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test bind validator", func(t *testing.T) {
		binding := Binding{
			BindingID: "",
		}
		err := binding.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("BindingID: cannot be blank"))
	})
}

func TestClient_ValidateExpiry(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test expiry is ok", func(t *testing.T) {
		binding := Binding{
			BindingID: "123",
			NewExpiry: "123123",
		}
		err := binding.Validate()
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_Bind(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Test Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		binding := Binding{
			BindingID: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
			NewExpiry: "123123",
		}
		testServer.Mux.HandleFunc(endpoints.Register, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		client := Client{
			API: acquiring.GetAPI(),
		}

		_, _, err := bind(context.Background(), client, endpoints.Register, binding)
		Expect(err).To(HaveOccurred())
	})
}

func TestClient_GetBindings(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Trigger GetBinding error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetBindings, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		_, _, err := GetBindings(context.Background(), "123123", nil)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test GetBinding with fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub
		_, _, err := GetBindings(context.Background(), "123123", nil)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("GetBinding is working as expected", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetBindings, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.OrderStatusResponse{
				OrderNumber: "70906e55-7114-41d6-8332-4609dc6590f4",
				ErrorCode:   0,
			})
		})

		_, _, err := GetBindings(context.Background(), "123123", nil)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_GetBindingsByCardOrId(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Trigger GetBindingsByCardOrId error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetBindingsByCardOrId, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		binding := GetBindingsRequest{
			UserName: "username",
			Password: "password",
		}

		_, _, err := GetBindingsByCardOrId(context.Background(), binding)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test GetBindingsByCardOrId with fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		binding := GetBindingsRequest{
			UserName: "username",
			Password: "password",
		}

		_, _, err := GetBindingsByCardOrId(context.Background(), binding)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("GetBindingByCardOrId is working as expected", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetBindingsByCardOrId, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsByCardOrIdResponse{
				ErrorCode:    0,
				ErrorMessage: "Binding is active",
			})
		})
		pan := "123123123123"
		binding := GetBindingsRequest{
			UserName: "username",
			Password: "password",
			PAN:      &pan,
		}
		_, _, err := GetBindingsByCardOrId(context.Background(), binding)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_CreateBindingNoPayment(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Trigger CreateBindingNoPayment error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.CreateBindingNoPayment, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		binding := CreateBindingNoPaymentRequest{
			UserName: "username",
			Password: "password",
		}

		_, _, err := CreateBindingNoPayment(context.Background(), binding)
		// We dont care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

	t.Run("Test CreateBindingNoPayment with fail on NewRestRequest", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		oldNewRequest := acquiring.NewRestRequest
		acquiring.NewRestRequest = NewRestRequestStub

		binding := CreateBindingNoPaymentRequest{
			UserName: "username",
			Password: "password",
		}

		_, _, err := CreateBindingNoPayment(context.Background(), binding)
		// We dont care what underlying error happened, we just don't run server to handle request
		Expect(err).To(HaveOccurred())
		acquiring.NewRestRequest = oldNewRequest
	})

	t.Run("CreateBindingNoPayment is working as expected", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.CreateBindingNoPayment, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsNoPaymentResponse{
				ErrorCode:    0,
				ErrorMessage: "Binding is active",
			})
		})
		binding := CreateBindingNoPaymentRequest{
			UserName:       "username",
			Password:       "password",
			PAN:            "123123123123",
			ClientId:       "123123123123",
			ExpiryDate:     "2021-01-01",
			CardHolderName: "username",
		}
		_, _, err := CreateBindingNoPayment(context.Background(), binding)
		// We dont care what underlying error happened
		Expect(err).ToNot(HaveOccurred())
	})
}
