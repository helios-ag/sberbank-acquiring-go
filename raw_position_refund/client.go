package raw_position_refund

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

// ProcessRawPositionRefundRequest Request is used to carry data related that passed to acquiring api requests.
type ProcessRawPositionRefundRequest struct {
	UserName            string
	Password            string
	Language            string
	OrderId             string
	Amount              int64
	PositionId          string
	AdditionalOfdParams *ofd.AdditionalOfdParams
}

func (request ProcessRawPositionRefundRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.OrderId, validation.Required, validation.Length(1, 30)),
		validation.Field(&request.Amount, validation.Required, validation.Min(1)),
		validation.Field(&request.UserName, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

// ProcessRawPositionRefund ProcessRawPositionRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:processrawpositionrefund
func ProcessRawPositionRefund(ctx context.Context, processRawPositionRefundRequest ProcessRawPositionRefundRequest) (*schema.ProcessRawRefundResponse, *http.Response, error) {
	return getClient().ProcessRawPositionRefund(ctx, processRawPositionRefundRequest)
}

// ProcessRawPositionRefund ProcessRawPositionRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:ProcessRawPositionRefund
func (c Client) ProcessRawPositionRefund(ctx context.Context, ProcessRawPositionRefundRequest ProcessRawPositionRefundRequest) (*schema.ProcessRawRefundResponse, *http.Response, error) {
	path := endpoints.ProcessRawPositionRefund

	if err := ProcessRawPositionRefundRequest.Validate(); err != nil {
		return nil, nil, err
	}
	var additionalOfdParams, _ = json.Marshal(ProcessRawPositionRefundRequest.AdditionalOfdParams)

	body := map[string]string{
		"userName": ProcessRawPositionRefundRequest.UserName,
		"password": ProcessRawPositionRefundRequest.Password,
		"amount":   strconv.FormatInt(ProcessRawPositionRefundRequest.Amount, 10),
	}
	body["language"] = ProcessRawPositionRefundRequest.Language
	body["orderId"] = ProcessRawPositionRefundRequest.OrderId
	body["amount"] = strconv.FormatInt(ProcessRawPositionRefundRequest.Amount, 10)
	body["positionId"] = ProcessRawPositionRefundRequest.PositionId

	if ProcessRawPositionRefundRequest.AdditionalOfdParams != nil {
		body["additionalOfdParams"] = string(additionalOfdParams)
	}

	var ProcessRawPositionRefundResponse schema.ProcessRawRefundResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodPost, path, body, nil)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &ProcessRawPositionRefundResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&ProcessRawPositionRefundResponse)

	return &ProcessRawPositionRefundResponse, result, err
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
