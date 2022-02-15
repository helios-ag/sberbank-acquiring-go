package bind

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"net/http"
	"regexp"
)

type Client struct {
	API acquiring.API
}

// Binding is used to make binding related requests
type Binding struct {
	bindingID  string
	newExpiry  string
	JSONParams map[string]string
}

func (binding Binding) Validate() error {
	return validation.ValidateStruct(&binding,
		validation.Field(&binding.bindingID, validation.Required),
		validation.Field(&binding.newExpiry, validation.Match(regexp.MustCompile("^[0-9]{6}$"))),
	)
}

// BindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:bindcard
func BindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	return getClient().BindCard(ctx, binding)
}

// BindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:bindcard
func (c Client) BindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.BindCard

	return bind(ctx, c, path, binding)
}

// UnBindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:unbindcard
func UnBindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	return getClient().UnBindCard(ctx, binding)
}

// UnBindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:unbindcard
func (c Client) UnBindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.UnBindCard

	return bind(ctx, c, path, binding)
}

var bind = func(ctx context.Context, client Client, path string, binding Binding) (*schema.Response, *http.Response, error) {
	if err := binding.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingID

	return client.bind(ctx, path, body, binding.JSONParams)
}

// ExtendBinding request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:extendbinding
func ExtendBinding(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	return getClient().ExtendBinding(ctx, binding)
}

// ExtendBinding request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:extendbinding
func (c Client) ExtendBinding(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.ExtendBinding

	if err := binding.Validate(); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingID
	body["newExpiry"] = binding.newExpiry

	return c.bind(ctx, path, body, binding.JSONParams)
}

func (c Client) bind(ctx context.Context, path string, body map[string]string, jsonParams map[string]string) (*schema.Response, *http.Response, error) {
	var response schema.Response
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, jsonParams)

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

// GetBindings request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindings
func GetBindings(ctx context.Context, clientId string, jsonParams map[string]string) (*schema.BindingsResponse, *http.Response, error) {
	return getClient().GetBindings(ctx, clientId, jsonParams)
}

// GetBindings request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindings
func (c Client) GetBindings(ctx context.Context, clientId string, jsonParams map[string]string) (*schema.BindingsResponse, *http.Response, error) {
	path := endpoints.GetBindings

	if len(clientId) > 255 {
		return nil, nil, fmt.Errorf("clientId is too long (max 255)")
	}

	body := make(map[string]string)
	body["clientId"] = clientId

	var response schema.BindingsResponse
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, jsonParams)

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
