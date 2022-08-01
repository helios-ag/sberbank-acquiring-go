package enrollment

import (
	"context"
	"encoding/json"
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

func TestClient_VerifyEnrollment(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Validate Enrollment PAN", func(t *testing.T) {
		prepareClient("localhost")

		enrollment := "4111111111111111111111111"

		_, _, err := VerifyEnrollment(context.Background(), enrollment)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("pan length shouldn't be less 13 or more 19 symbols"))
	})

	t.Run("Test Enrollment response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.VerifyEnrollment, func(w http.ResponseWriter, r *http.Request) {
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

		response, _, err := VerifyEnrollment(context.Background(), enrollment)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ErrorCode":    Equal(0),
			"ErrorMessage": Equal("Успешно"),
			"Enrolled":     Equal(byte('Y')),
		})))
	})

	t.Run("Trigger verifyEnrollment error on Do", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.VerifyEnrollment, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		_, _, err := VerifyEnrollment(context.Background(), "411111111111111")
		// We don't care what underlying error happened
		Expect(err).To(HaveOccurred())
	})

}
