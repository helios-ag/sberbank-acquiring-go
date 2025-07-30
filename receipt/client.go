package receipt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// StatusRequest ReceiptStatusRequest is used for building GetReceipt request
type StatusRequest struct {
	OrderId     string
	OrderNumber string
	UUID        string
	JsonParams  map[string]string
}

// GetReceiptStatus request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getreceiptstatus
func GetReceiptStatus(ctx context.Context, receiptStatusRequest StatusRequest) (*schema.ReceiptStatus, *http.Response, error) {
	return getClient().GetReceiptStatus(ctx, receiptStatusRequest)
}

// GetReceiptStatus request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getreceiptstatus
func (c Client) GetReceiptStatus(ctx context.Context, receiptStatusRequest StatusRequest) (*schema.ReceiptStatus, *http.Response, error) {
	path := endpoints.GetReceiptStatus

	if err := validateReceiptStatusRequest(receiptStatusRequest); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = receiptStatusRequest.OrderId
	body["orderNumber"] = receiptStatusRequest.OrderNumber
	body["uuid"] = receiptStatusRequest.UUID

	var response schema.ReceiptStatus

	req, err := c.API.NewRestRequest(ctx, "GET", path, body, receiptStatusRequest.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateReceiptStatusRequest(receiptStatusRequest StatusRequest) error {
	if receiptStatusRequest.OrderNumber == "" && receiptStatusRequest.OrderId == "" && receiptStatusRequest.UUID == "" {
		return fmt.Errorf("pass orderNumber or orderId or UUID")
	}

	if receiptStatusRequest.OrderNumber != "" {
		if len(receiptStatusRequest.OrderNumber) > 30 {
			return fmt.Errorf("orderNumber is too long (>30)")
		}
	}

	return nil
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
