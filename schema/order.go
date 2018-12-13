package schema

// Response is mapped response received from Sberbank API
type Response struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// OrderResponse is response received from RegisterOrder and RegisterOrderPreAuth requests
type OrderResponse struct {
	OrderId      string `json:"orderId,omitempty"`
	FormUrl      string `json:"formUrl,omitempty"`
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// OrderStatusResponse is response from GetOrderStatus request
type OrderStatusResponse struct {
	OrderNumber           string `json:"orderNumber"`
	OrderStatus           int    `json:"orderStatus,omitempty"`
	ActionCode            int    `json:"actionCode"`
	ActionCodeDescription string `json:"actionCodeDescription"`
	ErrorCode             int    `json:"errorCode,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	Amount                int    `json:"amount"`
	Currency              int    `json:"currency,omitempty"`
	Date                  string `json:"date"`
	OrderDescription      string `json:"orderDescription,omitempty"`
	Ip                    string `json:"ip"`
	MerchantOrderParams   []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"merchantOrderParams"`
	Attributes []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	CardAuthInfo struct {
		MaskedPan      string `json:"maskedPan,omitempty"`
		Expiration     int    `json:"expiration,omitempty"`
		CardholderName string `json:"cardholderName,omitempty"`
		ApprovalCode   string `json:"approvalCode,omitempty"`
		Chargeback     string `json:"chargeback,omitempty"`
		PaymentSystem  string `json:"paymentSystem"`
		Product        string `json:"product"`
		PaymentWay     string `json:"paymentWay"`
		SecureAuthInfo struct {
			Eci         int `json:"eci"`
			ThreeDSInfo struct {
				Xid string `json:"xid"`
			} `json:"threeDSInfo"`
		} `json:"secureAuthInfo"`
	} `json:"cardAuthInfo"`
	BankInfo struct {
		BankName        string `json:"bankName"`
		BankCountryCode string `json:"bankCountryCode"`
		BankCountryName string `json:"bankCountryName"`
	} `json:"bankName,omitempty"`
	TerminalId        string `json:"terminalId"`
	PaymentAmountInfo struct {
		ApprovedAmount  int    `json:"approvedAmount,omitempty"`
		DepositedAmount int    `json:"depositedAmount,omitempty"`
		RefundedAmount  int    `json:"refundedAmount,omitempty"`
		PaymentState    string `json:"paymentState"`
		FeeAmount       int    `json:"feeAmount"`
	} `json:"paymentAmountInfo,omitempty"`
}
