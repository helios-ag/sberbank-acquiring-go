package acquiring

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"net/http"
	"regexp"
	"strconv"
)

// Order is used to carry data related that passed to acquiring api requests.
//
// "OrderNumber" used to pass orderId to api
// "Amount" is amount of money (in pennies)
// "ReturnURL" in response
// "FailURL" in response
// "Description" check API Docs
// "PageView" custom pageview
// "MerchantLogin" check API Docs
// "ExpirationDate" check API Docs
// "BindingID" used in binding API
// "OrderBundle" OrderBundle data (cart to be consistent with 84 law and OFD 1.05)
// "AdditionalOfdParams" AdditionalOfdParams extra data (for OFD 1.05 and up)
// "Features" used in some endpoints of API
// "JSONParams" different json data that can be stored on api side
type Order struct {
	OrderNumber         string
	Amount              int
	ReturnURL           string
	FailURL             string
	Description         string
	PageView            string
	MerchantLogin       string
	ExpirationDate      string
	BindingID           string
	OrderBundle         OrderBundle
	AdditionalOfdParams AdditionalOfdParams
	Features            string
	JSONParams          map[string]string
}

type AdditionalOfdParams struct {
	AgentInfoType                   int      `json:"agent_info.type"`
	AgentInfoPayingOperation        string   `json:"agent_info.paying.operation,omitempty"`
	AgentInfoPayingPhones           []string `json:"agent_info.paying.phones,omitempty"`
	AgentInfoPaymentsOperatorPhones []string `json:"agent_info.paymentsOperator.phones,omitempty"`
	AgentInfoMTOperatorAddress      string   `json:"agent_info.MTOperator.address,omitempty"`
	AgentInfoMTOperatorInn          string   `json:"agent_info.MTOperator.inn,omitempty"`
	AgentInfoMTOperatorName         string   `json:"agent_info.MTOperator.name,omitempty"`
	AgentInfoMTOperatorPhones       []string `json:"agent_info.MTOperator.phones,omitempty"`
	SupplierInfoPhones              []string `json:"supplier_info.phones,omitempty"`
	Cashier                         string   `json:"cashier,omitempty"`
	AdditionalCheckProps            string   `json:"additional_check_props,omitempty"`
	AdditionalUserPropsName         string   `json:"additional_user_props.name,omitempty"`
	AdditionalUserPropsValue        string   `json:"additional_user_props.value,omitempty"`
	CashierInn                      string   `json:"cashier_inn,omitempty"`
	ClientAddress                   string   `json:"client.address,omitempty"`
	ClientBirthDate                 string   `json:"client.birth_date,omitempty"`
	ClientCitizenship               string   `json:"client.citizenship,omitempty"`
	ClientDocumentCode              string   `json:"client.document_code,omitempty"`
	ClientPassportNumber            string   `json:"client.passport_number,omitempty"`
	ClientMail                      string   `json:"client.email,omitempty"`
	ClientPhone                     string   `json:"client.phone,omitempty"`
	ClientInn                       string   `json:"client.inn,omitempty"`
	ClientName                      string   `json:"client.name,omitempty"`
	OperatingCheckPropsName         string   `json:"operatingCheckProps.name,omitempty"`
	OperatingCheckPropsTimestamp    string   `json:"operatingCheckProps.timestamp,omitempty"`
	OperatingCheckPropsValue        string   `json:"operatingCheckProps.value,omitempty"`
	SectoralCheckPropsDate          string   `json:"sectoralCheckProps.date,omitempty"`
	SectoralCheckPropsFederalId     string   `json:"sectoralCheckProps.federalId,omitempty"`
	SectoralCheckPropsNumber        string   `json:"sectoralCheckProps.number,omitempty"`
	SectoralCheckPropsValue         string   `json:"sectoralCheckProps.value,omitempty"`
}

type OrderBundle struct {
	OrderCreationDate string           `json:"orderCreationDate,omitempty"`
	CustomerDetails   *CustomerDetails `json:"customerDetails,omitempty"`
	CartItems         CartItems        `json:"cartItems"`
}

func (orderBundle OrderBundle) Validate() error {
	return validation.ValidateStruct(&orderBundle,
		validation.Field(&orderBundle.OrderCreationDate, validation.Length(1, 21)),
		validation.Field(&orderBundle.CustomerDetails),
	)
}

type CustomerDetails struct {
	Contact  string `json:"contact,omitempty"`
	Email    string `json:"email"`
	FullName string `json:"fullName,omitempty"`
	Passport string `json:"passport,omitempty"`
	Inn      string `json:"inn,omitempty"`

	Phone        string       `json:"phone,omitempty"`
	DeliveryInfo DeliveryInfo `json:"delivery_info,omitempty"`
}

func (customerDetails CustomerDetails) Validate() error {
	return validation.ValidateStruct(&customerDetails,
		validation.Field(&customerDetails.Contact, validation.Length(1, 40)),
		validation.Field(&customerDetails.Email, validation.Required, is.Email, validation.Length(1, 40)),
		validation.Field(&customerDetails.Phone, validation.Required, validation.Length(1, 40), validation.Match(regexp.MustCompile("^(('+7'|7|8)?([0-9]){10})$"))),
		validation.Field(&customerDetails.FullName, validation.Length(1, 100)),
		validation.Field(&customerDetails.Passport, validation.Length(1, 100)),
		validation.Field(&customerDetails.Inn, validation.RuneLength(1, 12)),
		validation.Field(&customerDetails.DeliveryInfo),
	)
}

type DeliveryInfo struct {
	DeliveryType    *string `json:"delivery_type,omitempty"`
	DeliveryCountry string  `json:"delivery_country"`
	DeliveryCity    string  `json:"delivery_city"`
	PostAddress     string  `json:"post_address"`
}

func (deliveryInfo DeliveryInfo) Validate() error {
	return validation.ValidateStruct(&deliveryInfo,
		validation.Field(&deliveryInfo.DeliveryType, validation.Length(1, 20)),
		validation.Field(&deliveryInfo.DeliveryCountry, validation.Required, validation.Length(1, 20)),
		validation.Field(&deliveryInfo.DeliveryCity, validation.Required, validation.Length(1, 40)),
		validation.Field(&deliveryInfo.PostAddress, validation.Required, validation.Length(1, 255)),
	)
}

type CartItems struct {
	Items []Item `json:"items"`
}

type Item struct {
	PositionId     string            `json:"positionId"`
	Name           string            `json:"name"`
	ItemDetails    ItemDetailsParams `json:"itemDetails,omitempty"`
	Quantity       Quantity          `json:"quantity"`
	ItemAmount     int               `json:"itemAmount,omitempty"`
	ItemCurrency   int               `json:"itemCurrency,omitempty"`
	ItemCode       string            `json:"itemCode"`
	ItemPrice      string            `json:"itemPrice"`
	ItemAttributes ItemAttributes    `json:"itemAttributes,omitempty"`
	Discount       Discount          `json:"discount,omitempty"`
	AgentInterest  AgentInterest     `json:"agentInterest,omitempty"`
	Tax            Tax               `json:"tax,omitempty"`
}

func (item Item) Validate() error {
	return validation.ValidateStruct(&item,
		validation.Field(&item.PositionId, validation.Length(1, 20)),
		validation.Field(&item.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&item.Quantity),
		validation.Field(&item.ItemCode, validation.Required, validation.Length(1, 100)),
		validation.Field(&item.ItemPrice, validation.Required, validation.Length(1, 18)),
	)
}

// Discount structure
type Discount struct {
	DiscountType  string `json:"discountType"`
	DiscountValue string `json:"discountValue"`
}

// Validate Validates Discount struct
func (discount Discount) Validate() error {
	return validation.ValidateStruct(&discount,
		validation.Field(&discount.DiscountType, validation.Length(1, 20)),
		validation.Field(&discount.DiscountValue, validation.Required, validation.Length(1, 20)),
	)
}

type AgentInterest struct {
	InterestType  string `json:"interestType"`
	InterestValue string `json:"interestValue"`
}

func (agentInterest AgentInterest) Validate() error {
	return validation.ValidateStruct(&agentInterest,
		validation.Field(&agentInterest.InterestType, validation.Required, validation.Length(1, 20)),
		validation.Field(&agentInterest.InterestValue, validation.Required, validation.Length(1, 20)),
	)
}

type Quantity struct {
	Value   int    `json:"value"`
	Measure string `json:"measure"`
}

func (quantity Quantity) Validate() error {
	return validation.ValidateStruct(&quantity,
		validation.Field(&quantity.Value, validation.Required, validation.Length(1, 18)),
		validation.Field(&quantity.Measure, validation.Required, validation.Length(1, 20)),
	)
}

type ItemDetails struct {
	ItemDetailsParams []ItemDetailsParams `json:"itemDetailsParams"`
}

func (itemDetails ItemDetails) Validate() error {
	return validation.ValidateStruct(&itemDetails,
		validation.Field(&itemDetails.ItemDetailsParams),
	)
}

type ItemAttributes struct {
	Attributes []Attributes `json:"attributes"`
}

func (itemAttributes ItemAttributes) Validate() error {
	return validation.ValidateStruct(&itemAttributes,
		validation.Field(&itemAttributes.Attributes),
	)
}

type Attributes struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (attributes Attributes) Validate() error {
	return validation.ValidateStruct(&attributes,
		validation.Field(&attributes.Name, validation.Required),
	)
}

type ItemDetailsParams struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

func (itemDetailsParams ItemDetailsParams) Validate() error {
	return validation.ValidateStruct(&itemDetailsParams,
		validation.Field(&itemDetailsParams.Value, validation.Required, validation.Length(1, 255)),
		validation.Field(&itemDetailsParams.Name, validation.Required, validation.Length(1, 255)),
	)
}

type Tax struct {
	TaxType int `json:"taxType,omitempty"`
	TaxSum  int `json:"taxSum,omitempty"`
}

func (tax Tax) Validate() error {
	return validation.ValidateStruct(&tax,
		validation.Field(&tax.TaxType, validation.Length(1, 2)),
		validation.Field(&tax.TaxSum, validation.Length(1, 18)),
	)
}

// RegisterOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:register
func (c *Client) RegisterOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Register

	if err := validateRegisterOrder(order); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

// RegisterOrderPreAuth request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:registerpreauth
func (c *Client) RegisterOrderPreAuth(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.RegisterPreAuth

	if err := validateRegisterOrder(order); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

func validateRegisterOrder(order Order) error {

	err := validation.ValidateStruct(&order,
		validation.Field(&order.OrderNumber, validation.Required, validation.Length(0, 30)),
		validation.Field(&order.ReturnURL, validation.Required, is.URL),
		validation.Field(&order.FailURL, validation.Required, is.URL),
		validation.Field(&order.FailURL, validation.Required, is.URL),
		validation.Field(&order.FailURL, validation.Required, is.URL),
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) register(ctx context.Context, path string, order Order) (*schema.OrderResponse, *http.Response, error) {
	body := make(map[string]string)
	var orderBundle, _ = json.Marshal(order.OrderBundle)
	body["orderNumber"] = order.OrderNumber
	body["amount"] = strconv.Itoa(order.Amount)
	body["returnUrl"] = order.ReturnURL
	body["failUrl"] = order.FailURL
	body["description"] = order.Description
	body["pageView"] = order.PageView
	body["merchantLogin"] = order.MerchantLogin
	body["expirationDate"] = order.ExpirationDate
	body["bindingId"] = order.BindingID
	body["orderBundle"] = string(orderBundle[:])
	body["features"] = order.Features

	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}

	var orderResponse schema.OrderResponse
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

	return &orderResponse, result, err
}

// Deposit request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:deposit
func (c *Client) Deposit(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Deposit

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["amount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}

	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

	return &orderResponse, result, err
}

// ReverseOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:reverse
func (c *Client) ReverseOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Reverse

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

	return &orderResponse, result, err
}

// RefundOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:refund
func (c *Client) RefundOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Refund

	if err := validateRefundOrder(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["refundAmount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

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

// GetOrderStatus request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getorderstatusextended
func (c *Client) GetOrderStatus(ctx context.Context, order Order) (*schema.OrderStatusResponse, *http.Response, error) {
	path := endpoints.GetOrderStatusExtended

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderStatusResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

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

// VerifyEnrollment request
// Checks if card enrolled in 3D Sec
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:verifyEnrollment
func (c *Client) VerifyEnrollment(ctx context.Context, pan string) (*schema.EnrollmentResponse, *http.Response, error) {
	path := endpoints.VerifyEnrollment

	if err := validatePan(pan); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["pan"] = pan

	var enrollmentResponse schema.EnrollmentResponse
	req, err := c.NewRestRequest(ctx, "GET", path, body, nil)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.Do(req, &enrollmentResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&enrollmentResponse)

	return &enrollmentResponse, result, err
}

func validatePan(pan string) error {
	panLength := len(pan)
	if panLength < 13 || panLength > 19 {
		return fmt.Errorf("pan length shouldn't be less 13 or more 19 symbols")
	}
	return nil
}

// UpdateSSLCardList request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:updateSSLCardList
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// Binding is used to make binding related requests
type Binding struct {
	bindingID  string
	newExpiry  int
	JSONParams map[string]string
}

// BindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:bindcard
func (c *Client) BindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.BindCard

	return bind(ctx, c, path, binding)
}

// UnBindCard request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:unbindcard
func (c *Client) UnBindCard(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.UnBindCard

	return bind(ctx, c, path, binding)
}

var bind = func(ctx context.Context, client *Client, path string, binding Binding) (*schema.Response, *http.Response, error) {
	if err := validateBind(binding); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingID

	return client.bind(ctx, path, body, binding.JSONParams)
}

// ExtendBinding request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:extendbinding
func (c *Client) ExtendBinding(ctx context.Context, binding Binding) (*schema.Response, *http.Response, error) {
	path := endpoints.ExtendBinding

	if err := validateBind(binding); err != nil {
		return nil, nil, err
	}

	if err := validateExpiry(binding); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["bindingId"] = binding.bindingID
	body["newExpiry"] = strconv.Itoa(binding.newExpiry)

	return c.bind(ctx, path, body, binding.JSONParams)
}

func validateBind(binding Binding) error {
	if binding.bindingID == "" {
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// GetBindings request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getbindings
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

// ReceiptStatusRequest is used for building GetReceipt request
type ReceiptStatusRequest struct {
	OrderId     string
	OrderNumber string
	UUID        string
	JsonParams  map[string]string
}

// GetReceiptStatus request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getreceiptstatus
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
	_ = json.NewDecoder(result.Body).Decode(&response)

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

// ApplePaymentRequest is used for building PayWithApplePay request
type ApplePaymentRequest struct {
	OrderNumber          string `json:"orderNumber"`
	Merchant             string `json:"merchant"`
	PaymentToken         string `json:"paymentToken"`
	Description          string `json:"description,omitempty"`
	PreAuth              bool   `json:"preAuth,omitempty"`
	AdditionalParameters map[string]string
}

// PayWithApplePay request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:payment_applepay
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateApplePaymentRequest(request ApplePaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}

// GooglePaymentRequest request
// "OrderNumber" _required_ used to pass orderId to api
// "Merchant" _required_ merchant is API
// "PaymentToken" _required_ token encoded in base64
// "Language" in response
// "Description" check API Docs
// "ClientId" client ID
// "IP" clients ip to store
// "PreAuth" preAuth flag
// "Amount" sum in penny
// "AdditionalParameters" used in some endpoints of API
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

// PayWithGooglePay request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:payment_googlepay
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateGooglePayRequest(request GooglePaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}

// SamsungPaymentRequest is used to send PayWithSamsungPay request
// "OrderNumber" _required_ used to pass orderId to api
// "Merchant" _required_ merchant is API
// "PaymentToken" _required_ token encoded in base64
// "Language" in response
// "Description" check API Docs
// "ClientId" client ID
// "IP" clients ip to store
// "PreAuth" preauth flag
// "AdditionalParameters" used in some endpoints of API
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

// PayWithSamsungPay is used to send PayWithSamsungPay request
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
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}

func validateSamsungPaymentRequest(request SamsungPaymentRequest) error {
	if request.OrderNumber == "" || request.Merchant == "" || request.PaymentToken == "" {
		return fmt.Errorf("orderNumber, merchant and PaymentToken are required")
	}
	return nil
}
