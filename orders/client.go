package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

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

func (order Order) Validate() error {
	return validation.ValidateStruct(&order,
		validation.Field(&order.ReturnURL, validation.Required, is.URL),
		validation.Field(&order.OrderNumber, validation.Required, validation.Length(1, 30)),
		validation.Field(&order.FailURL, is.URL),
	)
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
		validation.Field(&orderBundle.CartItems),
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
		validation.Field(&customerDetails.Phone, validation.Required, validation.Length(1, 12), validation.Match(regexp.MustCompile(`^((\\+7|7|8)+([0-9]){10})$`))),
		validation.Field(&customerDetails.FullName, validation.Length(1, 100)),
		validation.Field(&customerDetails.Passport, validation.Length(1, 100)),
		validation.Field(&customerDetails.Inn, validation.RuneLength(1, 12)),
		validation.Field(&customerDetails.DeliveryInfo),
	)
}

type DeliveryInfo struct {
	DeliveryType    string `json:"delivery_type,omitempty"`
	DeliveryCountry string `json:"delivery_country"`
	DeliveryCity    string `json:"delivery_city"`
	PostAddress     string `json:"post_address"`
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
func RegisterOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	return getClient().RegisterOrder(ctx, order)
}

// RegisterOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:register
func (c Client) RegisterOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Register

	if err := order.Validate(); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

// RegisterOrderPreAuth request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:registerpreauth
func RegisterOrderPreAuth(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	return getClient().RegisterOrderPreAuth(ctx, order)
}

// RegisterOrderPreAuth request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:registerpreauth
func (c Client) RegisterOrderPreAuth(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.RegisterPreAuth

	if err := order.Validate(); err != nil {
		return nil, nil, err
	}

	orderResponse, result, err := c.register(ctx, path, order)

	return orderResponse, result, err
}

//
func (c Client) register(ctx context.Context, path string, order Order) (*schema.OrderResponse, *http.Response, error) {
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

	req, err := c.API.NewRestRequest(ctx, http.MethodGet, path, body, order.JSONParams)

	if err != nil {
		return nil, nil, err
	}

	var orderResponse schema.OrderResponse
	result, err := c.API.Do(req, &orderResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&orderResponse)

	return &orderResponse, result, err
}

// Deposit request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:deposit
func Deposit(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	return getClient().Deposit(ctx, order)
}

// Deposit request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:deposit
func (c Client) Deposit(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Deposit

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["amount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

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

// ReverseOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:reverse
func ReverseOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	return getClient().ReverseOrder(ctx, order)
}

// ReverseOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:reverse
func (c Client) ReverseOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Reverse

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderResponse
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

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

// RefundOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:refund
func RefundOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	return getClient().RefundOrder(ctx, order)
}

// RefundOrder request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:refund
func (c Client) RefundOrder(ctx context.Context, order Order) (*schema.OrderResponse, *http.Response, error) {
	path := endpoints.Refund

	if err := validateRefundOrder(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber
	body["refundAmount"] = strconv.Itoa(order.Amount)

	var orderResponse schema.OrderResponse
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

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
func GetOrderStatus(ctx context.Context, order Order) (*schema.OrderStatusResponse, *http.Response, error) {
	return getClient().GetOrderStatus(ctx, order)
}

// GetOrderStatus request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:getorderstatusextended
func (c Client) GetOrderStatus(ctx context.Context, order Order) (*schema.OrderStatusResponse, *http.Response, error) {
	path := endpoints.GetOrderStatusExtended

	if err := validateOrderNumber(order); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["orderId"] = order.OrderNumber

	var orderResponse schema.OrderStatusResponse
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, order.JSONParams)

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

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
