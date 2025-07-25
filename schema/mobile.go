package schema

// ApplePaymentResponse is  response from PayWithApplePay request
type ApplePaymentResponse struct {
	Data struct {
		OrderID string `json:"orderId"`
	} `json:"data"`
	OrderStatus struct {
		ActionCode            int    `json:"actionCode"`
		ActionCodeDescription string `json:"actionCodeDescription"`
		Amount                int    `json:"amount"`
		Attributes            []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attributes"`
		AuthDateTime int    `json:"authDateTime"`
		AuthRefNum   string `json:"authRefNum"`
		BankInfo     struct {
			BankCountryName string `json:"bankCountryName"`
		} `json:"bankInfo"`
		CardAuthInfo struct {
			ApprovalCode   string `json:"approvalCode"`
			CardholderName string `json:"cardholderName"`
			Expiration     string `json:"expiration"`
			Pan            string `json:"pan"`
		} `json:"cardAuthInfo"`
		Currency            string `json:"currency"`
		Date                int    `json:"date"`
		ErrorCode           string `json:"errorCode"`
		IP                  string `json:"ip"`
		MerchantOrderParams []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"merchantOrderParams"`
		OrderNumber       string `json:"orderNumber"`
		OrderStatus       int    `json:"orderStatus"`
		PaymentAmountInfo struct {
			ApprovedAmount  int    `json:"approvedAmount"`
			DepositedAmount int    `json:"depositedAmount"`
			PaymentState    string `json:"paymentState"`
			RefundedAmount  int    `json:"refundedAmount"`
		} `json:"paymentAmountInfo"`
		TerminalID string `json:"terminalId"`
	} `json:"orderStatus"`
	Success bool `json:"success"`
	Error   struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Message     string `json:"message"`
	} `json:"error,omitempty"`
}

// SamsungPaymentResponse is response from PayWithSamsungPay request
type SamsungPaymentResponse struct {
	Success bool `json:"success"`
	Data    struct {
		OrderID string `json:"orderId,omitempty"`
	} `json:"data"`
	Error struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		Message     string `json:"message"`
	} `json:"error,omitempty"`
}

// MirPayPaymentResponse представляет ответ от API при вызове метода /payment/mirpay
// "ErrorCode" _обязательный_ код ошибки (0 — если успешно)
// "ErrorMessage" необязательное текстовое описание ошибки
// "OrderId" _обязательный_ уникальный идентификатор заказа в системе банка
// "FormUrl" _обязательный_ URL на платёжную форму для оплаты через Mir Pay
type MirPayPaymentResponse struct {
	ErrorCode    int    `json:"errorCode"`              // Код результата операции (0 — успех)
	ErrorMessage string `json:"errorMessage,omitempty"` // Текстовое описание ошибки (если есть)
	OrderId      string `json:"orderId"`                // Идентификатор заказа в системе банка
	FormUrl      string `json:"formUrl"`                // URL платёжной формы Mir Pay
	Data         struct {
		OrderID   string `json:"orderId,omitempty"`
		BindingId string `json:"BindingId,omitempty"`
	} `json:"data"`
	Error struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		Message     string `json:"message"`
	} `json:"error,omitempty"`
}

// GooglePaymentResponse is response from PayWithGooglePay request
type GooglePaymentResponse struct {
	Success bool `json:"success"`
	Data    struct {
		OrderID string `json:"orderId,omitempty"`
	} `json:"data"`
	Error struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		Message     string `json:"message"`
	} `json:"error,omitempty"`
}
