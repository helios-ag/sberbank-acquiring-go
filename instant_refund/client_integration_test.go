//go:build integration

package instant_refund_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/helios-ag/sberbank-acquiring-go/instant_refund"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	. "github.com/onsi/gomega"
)

// mockAPI is a mock implementation of the sberbank.API interface
type mockAPI struct {
	baseURL string
}

func (m *mockAPI) NewRestRequest(ctx context.Context, method, path string, data, queryParams map[string]string) (*http.Request, error) {
	// In a real implementation, you would create a proper HTTP request
	// For testing, we'll just create a simple request to our test server
	return http.NewRequest(method, m.baseURL+path, nil)
}

func (m *mockAPI) NewRequest(ctx context.Context, method, path string, data interface{}) (*http.Request, error) {
	// Not used in our tests, but required by the interface
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAPI) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// In a real implementation, you would make the HTTP request
	// For testing, we'll use the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into the provided interface
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return nil, err
		}
	}

	// Create a new response with the original body since we've already read it
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return resp, nil
}

// NewClient creates a new client for testing
func NewClient(serverURL string) *instant_refund.Client {
	return &instant_refund.Client{
		API: &mockAPI{baseURL: serverURL},
	}
}

func TestInstantRefund_Integration(t *testing.T) {
	RegisterTestingT(t)

	t.Run("when validation fails, returns error", func(t *testing.T) {
		// Given
		req := instant_refund.InstantRefundRequest{
			Amount: 1000,
		}

		// When
		_, _, err := instant_refund.InstantRefund(context.Background(), req)

		// Then
		Expect(err).To(HaveOccurred())
		// Check that the error contains the OrderNumber validation error
		Expect(err.Error()).To(ContainSubstring("OrderNumber: cannot be blank"))
	})
}

// TestInstantRefund_Unit tests the InstantRefund function with a mock HTTP server
func TestInstantRefund_Unit(t *testing.T) {
	RegisterTestingT(t)

	t.Run("when refund is successful, returns success response", func(t *testing.T) {
		// Create a test server that will respond with a success message
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.InstantRefundResponse{
				OrderID:      "test-order-123",
				ErrorCode:    0,
				ErrorMessage: "",
			})
		}))
		defer testServer.Close() // Clean up after the test

		// Create a client with the test server's URL
		client := NewClient(testServer.URL)

		// Given
		orderNum := "test-order-123"
		req := instant_refund.InstantRefundRequest{
			OrderNumber: &orderNum,
			Amount:      1000,
			UserName:    "test_user",
			Password:    "test_password",
		}

		// When
		resp, httpResp, err := client.InstantRefund(context.Background(), req)

		// Then
		Expect(err).NotTo(HaveOccurred())
		Expect(httpResp).NotTo(BeNil())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp).NotTo(BeNil())
		Expect(resp.OrderID).To(Equal("test-order-123"))
		Expect(resp.ErrorCode).To(Equal(0))
	})

	t.Run("when validation fails, returns error", func(t *testing.T) {
		// Given
		req := instant_refund.InstantRefundRequest{
			Amount: 1000,
		}

		// When
		_, _, err := instant_refund.InstantRefund(context.Background(), req)

		// Then
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("OrderNumber: cannot be blank"))
	})
}


