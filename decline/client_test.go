package bind

import (
	"context"
	"encoding/json"
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

func TestClient_Decline(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test decline Validate", func(t *testing.T) {
		decline := DeclineRequest{
			Username: "123",
			Password: "",
		}

		_, _, err := Decline(context.Background(), decline)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Password: cannot be blank."))
	})

	t.Run("Test Decline response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.Decline, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    0,
				ErrorMessage: "Success",
			})
		})

		decline := DeclineRequest{
			Username: "123",
			Password: "123123",
		}

		response, _, err := Decline(context.Background(), decline)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Success"),
		})))
	})

	t.Run("Test validate username with empty value", func(t *testing.T) {
		decline := DeclineRequest{
			Username: "",
			Password: "123123",
		}

		_, _, err := Decline(context.Background(), decline)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Username: cannot be blank."))
	})

	t.Run("Test Decline response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.Decline, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.BindingsResponse{
				ErrorCode:    0,
				ErrorMessage: "Success",
			})
		})

		declineRequest := DeclineRequest{
			Username: "user",
			Password: "password",
		}

		response, _, err := Decline(context.Background(), declineRequest)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Success"),
		})))
	})

}

func TestClient_ValidateBind(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test bind validator", func(t *testing.T) {
		decline := DeclineRequest{
			Username: "123",
		}
		err := decline.Validate()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Password: cannot be blank."))
	})
}

func TestClient_Decline_Do(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Test Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		binding := DeclineRequest{
			Username: "fd3afc57-c6d0-4e08-aaef-1b7cfeb093dc",
			Password: "123123",
		}
		testServer.Mux.HandleFunc(endpoints.Decline, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		client := Client{
			API: acquiring.GetAPI(),
		}

		_, _, err := decline(context.Background(), client, endpoints.Decline, binding)
		Expect(err).To(HaveOccurred())
	})
}
