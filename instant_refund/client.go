package instant_refund

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// InstantRefundRequest Order is used to carry data related that passed to acquiring api requests.
type InstantRefundRequest struct {
	UserName       string
	Password       string
	Amount         int64
	Language       *string
	Currency       *int
	OrderNumber    *string
	BindingID      *string
	SeToken        *string
	PAN            *string
	CVC            *string
	Expiry         *string
	CardHolderName *string
	JSONParams     *JSONParams `json:"jsonParams,omitempty"`
}

type JSONParams struct {
	CashboxID *int64 `json:"cashboxId,omitempty"` // Идентификатор кассы
	BasketID  *int64 `json:"basketId,omitempty"`  // Идентификатор корзины покупки или возврата
}

func (request InstantRefundRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.OrderNumber, validation.Required, validation.Length(1, 30)),
		validation.Field(&request.Amount, validation.Required, validation.Min(1)),
		validation.Field(&request.UserName, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

// InstantRefund InstantRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:instantrefund
func InstantRefund(ctx context.Context, instantRefundRequest InstantRefundRequest) (*schema.InstantRefundResponse, *http.Response, error) {
	return getClient().InstantRefund(ctx, instantRefundRequest)
}

// InstantRefund InstantRefundRequest request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:instantrefund
func (c Client) InstantRefund(ctx context.Context, instantRefundRequest InstantRefundRequest) (*schema.InstantRefundResponse, *http.Response, error) {
	path := endpoints.InstantRefund

	if err := instantRefundRequest.Validate(); err != nil {
		return nil, nil, err
	}
	var jsonParams, _ = json.Marshal(instantRefundRequest.JSONParams)

	body := map[string]string{
		"userName": instantRefundRequest.UserName,
		"password": instantRefundRequest.Password,
		"amount":   strconv.FormatInt(instantRefundRequest.Amount, 10),
	}
	if instantRefundRequest.Language != nil {
		body["language"] = *instantRefundRequest.Language
	}

	if instantRefundRequest.Currency != nil {
		body["currency"] = strconv.Itoa(*instantRefundRequest.Currency)
	}

	if instantRefundRequest.OrderNumber != nil {
		body["orderNumber"] = *instantRefundRequest.OrderNumber
	}

	if instantRefundRequest.BindingID != nil {
		body["bindingId"] = *instantRefundRequest.BindingID
	}

	if instantRefundRequest.SeToken != nil {
		body["seToken"] = *instantRefundRequest.SeToken
	}

	if instantRefundRequest.PAN != nil {
		body["pan"] = *instantRefundRequest.PAN
	}
	if instantRefundRequest.CVC != nil {
		body["cvc"] = *instantRefundRequest.CVC
	}

	if instantRefundRequest.Expiry != nil {
		body["expire"] = *instantRefundRequest.Expiry
	}

	if instantRefundRequest.CardHolderName != nil {
		body["cardHolderName"] = *instantRefundRequest.CardHolderName
	}

	if instantRefundRequest.JSONParams != nil {
		body["jsonParams"] = string(jsonParams)
	}

	var orderResponse schema.InstantRefundResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodPost, path, body, nil)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

	return &orderResponse, result, err
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
