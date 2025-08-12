package external_receipt

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "testing"

    "github.com/helios-ag/sberbank-acquiring-go/endpoints"
    "github.com/helios-ag/sberbank-acquiring-go/schema"
    server "github.com/helios-ag/sberbank-acquiring-go/testing"
    . "github.com/onsi/gomega"
)

func TestExternalReceipt_Additional(t *testing.T) {
    RegisterTestingT(t)

    t.Run("Test validation error - missing order ID", func(t *testing.T) {
        req := ExternalReceiptRequest{
            // Missing required fields
            UserName: "test",
            Password: "test",
        }

        _, _, err := GetExternalReceipt(context.Background(), req)
        Expect(err).To(HaveOccurred())
        Expect(err.Error()).To(ContainSubstring("userName and Password and mdOrder and Receipt are required"))
    })

    t.Run("Test successful request with order ID", func(t *testing.T) {
        testServer := server.NewServer()
        defer testServer.Teardown()
        prepareClient(testServer.URL)

        expectedResponse := schema.ExternalReceipt{
            ErrorCode:    0,
            ErrorMessage: "",
        }

        testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(expectedResponse)
        })

        req := ExternalReceiptRequest{
            MdOrder:  "order-123",
            UserName: "test",
            Password: "test",
            Receipt: &Receipt{
                PaymentType: 1,
            },
        }

        resp, httpResp, err := GetExternalReceipt(context.Background(), req)

        Expect(err).NotTo(HaveOccurred())
        Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
        Expect(resp.ErrorCode).To(Equal(0))
    })

    t.Run("Test successful request with order number", func(t *testing.T) {
        testServer := server.NewServer()
        defer testServer.Teardown()
        prepareClient(testServer.URL)

        expectedResponse := schema.ExternalReceipt{
            ErrorCode:    0,
            ErrorMessage: "",
        }

        testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(expectedResponse)
        })

        // API requires mdOrder; there is no OrderNumber field in request struct
        req := ExternalReceiptRequest{
            MdOrder:  "ORDER-456",
            UserName: "test",
            Password: "test",
            Receipt: &Receipt{PaymentType: 1},
        }

        resp, httpResp, err := GetExternalReceipt(context.Background(), req)

        Expect(err).NotTo(HaveOccurred())
        Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
        Expect(resp.ErrorCode).To(Equal(0))
    })

    t.Run("Test API error response", func(t *testing.T) {
        testServer := server.NewServer()
        defer testServer.Teardown()
        prepareClient(testServer.URL)

        errorResponse := map[string]any{
            "errorCode":    6,
            "errorMessage": "Receipt not found",
        }

        testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(errorResponse)
        })

        req := ExternalReceiptRequest{
            MdOrder:  "non-existent-order",
            UserName: "test",
            Password: "test",
            Receipt: &Receipt{PaymentType: 1},
        }

        _, httpResp, err := GetExternalReceipt(context.Background(), req)

        Expect(err).To(HaveOccurred())
        Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
    })

    t.Run("Test with additional parameters", func(t *testing.T) {
        testServer := server.NewServer()
        defer testServer.Teardown()
        prepareClient(testServer.URL)

        successResponse := schema.ExternalReceipt{
            ErrorCode:    0,
            ErrorMessage: "",
        }

        var receivedParams map[string]string
        testServer.Mux.HandleFunc(endpoints.ExternalReceipt, func(w http.ResponseWriter, r *http.Request) {
            // Read JSON body sent by NewRequest
            bodyBytes, _ := io.ReadAll(r.Body)
            _ = r.Body.Close()
            receivedParams = make(map[string]string)
            _ = json.Unmarshal(bodyBytes, &receivedParams)

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(successResponse)
        })

        cashboxID := int64(123)
        basketID := int64(456)
        req := ExternalReceiptRequest{
            MdOrder:  "order-123",
            UserName: "test",
            Password: "test",
            Receipt: &Receipt{
                PaymentType: 1,
            },
            JSONParams: &JSONParams{
                CashboxID: &cashboxID,
                BasketID:  &basketID,
            },
        }

        _, _, err := GetExternalReceipt(context.Background(), req)

        Expect(err).NotTo(HaveOccurred())
        Expect(receivedParams).To(HaveKeyWithValue("mdOrder", "order-123"))
        Expect(receivedParams).To(HaveKey("jsonParams"))
    })
}
