package acquiring

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"net/http"
	"net/url"
	"strconv"
)

type Order struct {
	OrderNumber    string
	Amount         int
	ReturnUrl      string
	FailUrl        string
	Description    string
	PageView       string
	MerchantLogin  string
	ExpirationDate string
	BindingId      string
	Features       string
	JsonParams     map[string]string
}

func (c *Client) RegisterOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Register

	if err := validateRegisterOrder(order); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

func (c *Client) RegisterOrderPreAuth(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.RegisterPreAuth

	if err := validateRegisterOrder(order); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

func validateRegisterOrder(order Order) error {
	if order.ReturnUrl == "" {
		return fmt.Errorf("return url must be set")
	}

	if _, err := url.Parse(order.ReturnUrl); err != nil {
		return fmt.Errorf("unable to parse ReturnUrl: %v", err)
	}

	if order.OrderNumber != "" {
		if len(order.OrderNumber) > 30 {
			return fmt.Errorf("orderNumber is too long (>30)")
		}
	}

	if order.FailUrl != "" {
		if _, err := url.Parse(order.FailUrl); err != nil {
			return fmt.Errorf("unable to parse FailUrl: %v", err)
		}
	}

	return nil
}

func (c *Client) register(ctx context.Context, path string, order Order) (*schema.OrderResponse, *http.Response, error) {
	body := make(map[string]string)
	body["orderNumber"] = order.OrderNumber
	body["amount"] = strconv.Itoa(order.Amount)
	body["returnUrl"] = order.ReturnUrl
	body["failUrl"] = order.FailUrl
	body["description"] = order.Description
	body["pageView"] = order.PageView
	body["merchantLogin"] = order.MerchantLogin
	body["expirationDate"] = order.ExpirationDate
	body["bindingId"] = order.BindingId
	body["features"] = order.Features
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	var orderResponse schema.OrderResponse
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&orderResponse)
	if err != nil {
		return nil, result, err
	}

	return &orderResponse, result, err
}

func (c *Client) Deposit(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Deposit

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["amount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&orderResponse)
	if err != nil {
		return nil, result, err
	}
	return &orderResponse, result, err
}

func (c *Client) ReverseOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Reverse

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&orderResponse)
	if err != nil {
		return nil, result, err
	}
	return &orderResponse, result, err
}

func (c *Client) RefundOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Refund

	if err := validateRefundOrder(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["refundAmount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&orderResponse)
	if err != nil {
		return nil, result, err
	}
	return &orderResponse, result, err
}

func validateRefundOrder(order Order) error {
	if order.OrderNumber == "" {
		return fmt.Errorf("orderNumber cant be empty")
	}

	if order.OrderNumber != "" {
		if len(order.OrderNumber) > 30 {
			return fmt.Errorf("orderNumber is too long (>30)")
		}
	}

	if order.Amount <= 0 {
		return fmt.Errorf("refund amount should be more 0")
	}

	return nil
}

func (c *Client) GetOrderStatus(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.GetOrderStatusExtended

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&orderResponse)
	if err != nil {
		return nil, result, err
	}
	return &orderResponse, result, err
}

func validateOrderNumber(order Order) error {
	if order.OrderNumber == "" {
		return fmt.Errorf("orderNumber cant be empty")
	}

	if order.OrderNumber != "" {
		if len(order.OrderNumber) > 30 {
			return fmt.Errorf("orderNumber is too long (>30)")
		}
	}

	return nil
}

type Enrollment struct {
	pan        string
	jsonParams map[string]string
}

func (c *Client) VerifyEnrollment(ctx context.Context, enrollment Enrollment) (*schema.EnrollmentResponse, *http.Response, error) {
	path := endpoints.VerifyEnrollment

	if err := validatePan(enrollment); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["pan"] = enrollment.pan

	var enrollmentResponse schema.EnrollmentResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, enrollment.jsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &enrollmentResponse)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&enrollmentResponse)
	if err != nil {
		return nil, result, err
	}
	return &enrollmentResponse, result, err
}

func validatePan(enrollment Enrollment) error {
	panLength := len(enrollment.pan)
	if panLength < 13 || panLength > 19 {
		return fmt.Errorf("pan length shouldn't be less 13 or more 19 symbols")
	}
	return nil
}

func (c *Client) UpdateSSLCardList(ctx context.Context, mdorder string, jsonParams map[string]string) (*schema.Response, *http.Response, error) {
	path := endpoints.UpdateSSLCardList
	body := make(map[string]string)
	body["mdorder"] = mdorder

	var response schema.Response
	req, err := c.NewRestRequest(ctx, "GET", path, body, jsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

type Binding struct {
	bindingId  string
	newExpiry  int
	jsonParams map[string]string
}

func (c *Client) BindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.BindCard

	if err := validateBind(binding); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingId

	return c.bind(ctx, path, body, binding.jsonParams)
}

func (c *Client) UnBindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.UnBindCard

	if err := validateBind(binding); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingId

	return c.bind(ctx, path, body, binding.jsonParams)
}

func (c *Client) ExtendBinding(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.ExtendBinding

	if err := validateBind(binding); err != nil {
		return nil, nil, err
	}

	if err := validateExpiry(binding); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingId
	body["newExpiry"] = strconv.Itoa(binding.newExpiry)

	return c.bind(ctx, path, body, binding.jsonParams)
}

func validateBind(binding Binding) error {
	if binding.bindingId == "" {
		return fmt.Errorf("bindingId can't be empty")
	}

	return nil
}

func validateExpiry(binding Binding) error {
	if len(strconv.Itoa(binding.newExpiry)) != 6 {
		return fmt.Errorf("new expiry date should have 6 digits")
	}

	return nil
}

func (c *Client) bind(ctx context.Context, path string, body map[string]string, jsonParams map[string]string) (*schema.Response, *http.Response, error) {
	var response schema.Response
	req, err := c.NewRestRequest(ctx, "GET", path, body, jsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

func (c *Client) GetBindings(ctx context.Context, clientId string, jsonParams map[string]string) (*schema.BindingsResponse, *http.Response, error) {
	path := endpoints.GetBindings

	if len(clientId) > 255 {
		return nil, nil, fmt.Errorf("clientId is too long (max 255)")
	}

	body := make(map[string]string)
	body["clientId"] = clientId

	var response schema.BindingsResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, jsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

type ReceiptStatusRequest struct {
	OrderId     string
	OrderNumber string
	UUID        string
	JsonParams  map[string]string
}

func (c *Client) GetReceiptStatus(ctx context.Context, receiptStatusRequest ReceiptStatusRequest) (*schema.ReceiptStatus, *http.Response, error) {
	path := endpoints.GetReceiptStatus

	if err := validateReceiptStatusRequest(receiptStatusRequest); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = receiptStatusRequest.OrderId
	body["orderNumber"] = receiptStatusRequest.OrderNumber
	body["uuid"] = receiptStatusRequest.UUID

	var response schema.ReceiptStatus

	req, err := c.NewRestRequest(ctx, "GET", path, body, receiptStatusRequest.JsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

func validateReceiptStatusRequest(receiptStatusRequest ReceiptStatusRequest) error {
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

type ApplePaymentRequest struct {
	OrderNumber          string `json:"orderNumber"`
	Merchant             string `json:"merchant"`
	PaymentToken         string `json:"paymentToken"`
	Description          string `json:"description,omitempty"`
	PreAuth              bool   `json:"preAuth,omitempty"`
	AdditionalParameters map[string]string
}

func (c *Client) PayWithApplePay(ctx context.Context, applePaymentRequest ApplePaymentRequest) (*schema.ApplePaymentResponse, *http.Response, error) {
	path := endpoints.ApplePay

	if err := validateApplePaymentRequest(applePaymentRequest); err != nil {
		return nil, nil, err
	}

	var response schema.ApplePaymentResponse

	req, err := c.NewRequest(ctx, "GET", path, applePaymentRequest)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

func validateApplePaymentRequest(request ApplePaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}

type GooglePaymentRequest struct {
	OrderNumber          string            `json:"orderNumber"`
	Merchant             string            `json:"merchant"`
	PaymentToken         string            `json:"paymentToken"`
	Language             string            `json:"language"`
	Description          string            `json:"description,omitempty"`
	ClientId             string            `json:"clientId,omitempty"`
	IP                   string            `json:"ip,omitempty"`
	Amount               int               `json:"amount,omitempty"`
	CurrencyCode         int               `json:"currencyCode,omitempty"`
	Email                string            `json:"email,omitempty"`
	Phone                string            `json:"phone,omitempty"`
	ReturnUrl            string            `json:"returnUrl"`
	FailUrl              string            `json:"failUrl,omitempty"`
	PreAuth              bool              `json:"preAuth,omitempty"`
	AdditionalParameters map[string]string `json:"additionalParameters"`
}

func (c *Client) PayWithGooglePay(ctx context.Context, googlePaymentRequest GooglePaymentRequest) (*schema.GooglePaymentResponse, *http.Response, error) {
	path := endpoints.GooglePay

	if err := validateGooglePayRequest(googlePaymentRequest); err != nil {
		return nil, nil, err
	}

	var response schema.GooglePaymentResponse

	req, err := c.NewRequest(ctx, "GET", path, googlePaymentRequest)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

func validateGooglePayRequest(request GooglePaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}

type SamsungPaymentRequest struct {
	OrderNumber          string            `json:"orderNumber"`
	Merchant             string            `json:"merchant"`
	PaymentToken         string            `json:"paymentToken"`
	Language             string            `json:"language"`
	Description          string            `json:"description,omitempty"`
	ClientId             string            `json:"clientId,omitempty"`
	IP                   string            `json:"ip,omitempty"`
	PreAuth              bool              `json:"preAuth,omitempty"`
	AdditionalParameters map[string]string `json:"additionalParameters"`
}

func (c *Client) PayWithSamsungPay(ctx context.Context, samsungPaymentRequest SamsungPaymentRequest) (*schema.SamsungPaymentResponse, *http.Response, error) {
	path := endpoints.SamsungPay

	if err := validateSamsungPaymentRequest(samsungPaymentRequest); err != nil {
		return nil, nil, err
	}

	var response schema.SamsungPaymentResponse
	req, err := c.NewRequest(ctx, "GET", path, samsungPaymentRequest)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		return nil, result, err
	}
	return &response, result, err
}

func validateSamsungPaymentRequest(request SamsungPaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}
