package endpoints

// API endpoints
const (
	Register                 string = "/payment/rest/register.do"
	RegisterPreAuth          string = "/payment/rest/registerPreAuth.do"
	Deposit                  string = "/payment/rest/deposit.do"
	Reverse                  string = "/payment/rest/reverse.do"
	Refund                   string = "/payment/rest/refund.do"
	InstantRefund            string = "/payment/rest/instantRefund.do"
	ProcessRawSumRefund      string = "/payment/rest/processRawSumRefund.do"
	ProcessRawPositionRefund string = "/payment/rest/processRawPositionRefund.do"

	GetOrderStatusExtended string = "/payment/rest/getOrderStatusExtended.do"
	GetReceiptStatus       string = "/payment/rest/getReceiptStatus.do"

	UnBindCard             string = "/payment/rest/unBindCard.do"
	BindCard               string = "/payment/rest/bindCard.do"
	GetBindings            string = "/payment/rest/getBindings.do"
	GetBindingsByCardOrId  string = "/payment/rest/getBindingsByCardOrId.do"
	ExtendBinding          string = "/payment/rest/extendBinding.do"
	CreateBindingNoPayment string = "/payment/rest/createBindingNoPayment.do"

	ApplePay        string = "/payment/applepay/payment.do"
	SamsungPay      string = "/payment/samsung/payment.do"
	SamsungWebPay   string = "/payment/samsungWeb/payment.do"
	GooglePay       string = "/payment/google/payment.do"
	MirPay          string = "/payment/mir/payment.do"
	MirPayDirect    string = "/payment/mir/paymentDirect.do"
	Decline         string = "/payment/rest/decline.do"
	ExternalReceipt string = "/fes-nspk-proxy/externalReceipt"

	VerifyEnrollment  string = "/payment/rest/verifyEnrollment.do"
	UpdateSSLCardList string = "/payment/rest/updateSSLCardList.do"
)
