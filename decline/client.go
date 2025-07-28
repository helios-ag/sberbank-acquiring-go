package bind

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// DeclineRequest is used to make Decline method related requests
type DeclineRequest struct {
	Username      string
	Password      string
	MerchantLogin string
	Language      string
	OrderNumber   string
	OrderId       string
}

func (decline DeclineRequest) Validate() error {
	return validation.ValidateStruct(&decline,
		validation.Field(&decline.Username, validation.Required),
		validation.Field(&decline.Password, validation.Required),
	)
}

// Decline request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:decline
func Decline(ctx context.Context, decline DeclineRequest) (*schema.DeclineResponse, *http.Response, error) {
	return getClient().Decline(ctx, decline)
}

// Decline request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:decline
func (c Client) Decline(ctx context.Context, declineRequest DeclineRequest) (*schema.DeclineResponse, *http.Response, error) {
	path := endpoints.Decline

	return decline(ctx, c, path, declineRequest)
}

var decline = func(ctx context.Context, client Client, path string, decline DeclineRequest) (*schema.DeclineResponse, *http.Response, error) {
	if err := decline.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["userName"] = decline.Username
	body["password"] = decline.Password
	body["merchantLogin"] = decline.MerchantLogin
	body["language"] = decline.Language
	body["orderId"] = decline.OrderId
	body["orderNumber"] = decline.OrderNumber

	return client.decline(ctx, path, body)
}

func (c Client) decline(ctx context.Context, path string, body map[string]string) (*schema.DeclineResponse, *http.Response, error) {
	var response schema.DeclineResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodGet, path, body, nil)

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

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
