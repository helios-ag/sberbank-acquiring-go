package bind

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// Binding is used to make binding related requests
type Binding struct {
	BindingID  string
	NewExpiry  string
	JSONParams map[string]string
}

func (binding Binding) Validate() error {
	return validation.ValidateStruct(&binding,
		validation.Field(&binding.BindingID, validation.Required),
		validation.Field(&binding.NewExpiry, validation.Match(regexp.MustCompile("^[0-9]{6}$"))),
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
	body["bindingId"] = binding.BindingID

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
	body["bindingId"] = binding.BindingID
	body["newExpiry"] = binding.NewExpiry

	return c.bind(ctx, path, body, binding.JSONParams)
}

func (c Client) bind(ctx context.Context, path string, body map[string]string, jsonParams map[string]string) (*schema.Response, *http.Response, error) {
	var response schema.Response
	req, err := c.API.NewRestRequest(ctx, http.MethodGet, path, body, jsonParams)

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

type GetBindingsRequest struct {
	PAN         *string
	UserName    string
	Password    string
	BindingID   *string
	ShowExpired *bool
}

func validateGetBindingRequest(request GetBindingsRequest) error {
	if request.UserName == "" || request.Password == "" {
		return fmt.Errorf("userName and Password are required")
	}
	if request.PAN == nil && request.BindingID == nil {
		return fmt.Errorf("either PAN or BindingID must be provided")
	}

	if len(*request.PAN) < 12 || len(*request.PAN) > 19 {
		return fmt.Errorf("PAN must be between 12 and 19 characters")
	}

	return nil
}

// GetBindingsByCardOrId request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindingsbycardorid
func GetBindingsByCardOrId(ctx context.Context, request GetBindingsRequest) (*schema.BindingsByCardOrIdResponse, *http.Response, error) {
	return getClient().GetBindingsByCardOrId(ctx, request)
}

// GetBindingsByCardOrId request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindingsbycardorid
func (c Client) GetBindingsByCardOrId(ctx context.Context, request GetBindingsRequest) (*schema.BindingsByCardOrIdResponse, *http.Response, error) {
	path := endpoints.GetBindingsByCardOrId

	if err := validateGetBindingRequest(request); err != nil {
		return nil, nil, err
	}

	body := map[string]string{
		"userName": request.UserName,
		"password": request.Password,
	}

	if request.PAN != nil {
		body["pan"] = *request.PAN
	}
	if request.BindingID != nil {
		body["bindingId"] = *request.BindingID
	}
	if request.ShowExpired != nil {
		body["showExpired"] = strconv.FormatBool(*request.ShowExpired)
	}

	var response schema.BindingsByCardOrIdResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodPost, path, body, nil)

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

type CreateBindingNoPaymentRequest struct {
	UserName             string
	Password             string
	ClientId             string
	CardHolderName       string
	PAN                  string
	MerchantLogin        *string
	ExpiryDate           string
	AdditionalParameters map[string]string
}

func validateCreateBindingNoPaymentRequest(request CreateBindingNoPaymentRequest) error {
	if request.UserName == "" || request.Password == "" || request.ClientId == "" || request.CardHolderName == "" || request.PAN == "" || request.ExpiryDate == "" {
		return fmt.Errorf("userName, Password, ClientId, CardHolderName, PAN, ExpiryDate are required")
	}

	return nil
}

// CreateBindingNoPayment request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindingsbycardorid
func CreateBindingNoPayment(ctx context.Context, request CreateBindingNoPaymentRequest) (*schema.BindingsNoPaymentResponse, *http.Response, error) {
	return getClient().CreateBindingNoPayment(ctx, request)
}

// CreateBindingNoPayment request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindingsbycardorid
func (c Client) CreateBindingNoPayment(ctx context.Context, request CreateBindingNoPaymentRequest) (*schema.BindingsNoPaymentResponse, *http.Response, error) {
	path := endpoints.CreateBindingNoPayment

	if err := validateCreateBindingNoPaymentRequest(request); err != nil {
		return nil, nil, err
	}

	body := map[string]string{
		"userName":       request.UserName,
		"password":       request.Password,
		"clientId":       request.ClientId,
		"cardHolderName": request.CardHolderName,
		"pan":            request.PAN,
		"expiryDate":     request.ExpiryDate,
	}

	if request.AdditionalParameters != nil {
		var params, _ = json.Marshal(request.AdditionalParameters)

		body["additionalParameters"] = string(params)
	}
	if request.MerchantLogin != nil {
		body["merchantLogin"] = *request.MerchantLogin
	}

	var response schema.BindingsNoPaymentResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodPost, path, body, nil)

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
