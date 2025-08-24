package raw_sum_refund

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/ofd"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// ProcessRawSumRefundRequest Request is used to carry data related that passed to acquiring api requests.
type ProcessRawSumRefundRequest struct {
	UserName            string
	Password            string
	Language            string
	OrderId             string
	Amount              int64
	Name                string
	JSONParams          *JSONParams `json:"jsonParams,omitempty"`
	AdditionalOfdParams *ofd.AdditionalOfdParams
	ItemCode            string `json:"itemCode"`
	TaxType             int    `json:"taxType,omitempty"`
}

type JSONParams struct {
	CashboxID *int64 `json:"cashboxId,omitempty"` // Идентификатор кассы
	BasketID  *int64 `json:"basketId,omitempty"`  // Идентификатор корзины покупки или возврата
}

func (request ProcessRawSumRefundRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.OrderId, validation.Required, validation.Length(1, 30)),
		validation.Field(&request.Amount, validation.Required, validation.Min(1)),
		validation.Field(&request.UserName, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

// ProcessRawSumRefund ProcessRawSumRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:processrawsumrefund
func ProcessRawSumRefund(ctx context.Context, processRawSumRefundRequest ProcessRawSumRefundRequest) (*schema.ProcessRawSumRefundResponse, *http.Response, error) {
	return getClient().ProcessRawSumRefund(ctx, processRawSumRefundRequest)
}

// ProcessRawSumRefund ProcessRawSumRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:processrawsumrefund
func (c Client) ProcessRawSumRefund(ctx context.Context, processRawSumRefundRequest ProcessRawSumRefundRequest) (*schema.ProcessRawSumRefundResponse, *http.Response, error) {
	path := endpoints.ProcessRawSumRefund

	if err := processRawSumRefundRequest.Validate(); err != nil {
		return nil, nil, err
	}
	var jsonParams, _ = json.Marshal(processRawSumRefundRequest.JSONParams)
	var additionalOfdParams, _ = json.Marshal(processRawSumRefundRequest.AdditionalOfdParams)

	body := map[string]string{
		"userName": processRawSumRefundRequest.UserName,
		"password": processRawSumRefundRequest.Password,
		"amount":   strconv.FormatInt(processRawSumRefundRequest.Amount, 10),
	}
	body["language"] = processRawSumRefundRequest.Language
	body["orderId"] = processRawSumRefundRequest.OrderId
	body["amount"] = strconv.FormatInt(processRawSumRefundRequest.Amount, 10)
	body["name"] = processRawSumRefundRequest.Name
	body["itemCode"] = processRawSumRefundRequest.ItemCode
	body["taxType"] = strconv.Itoa(processRawSumRefundRequest.TaxType)

	if processRawSumRefundRequest.JSONParams != nil {
		body["jsonParams"] = string(jsonParams)
	}
	if processRawSumRefundRequest.AdditionalOfdParams != nil {
		body["additionalOfdParams"] = string(additionalOfdParams)
	}

	var processRawSumRefundResponse schema.ProcessRawSumRefundResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodPost, path, body, nil)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &processRawSumRefundResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&processRawSumRefundResponse)

	return &processRawSumRefundResponse, result, err
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
