package endpoints

// API endpoints
const (
	Register               string = "/payment/rest/register.do"
	RegisterPreAuth        string = "/payment/rest/registerPreAuth.do"
	Deposit                string = "/payment/rest/deposit.do"
	Reverse                string = "/payment/rest/reverse.do"
	Refund                 string = "/payment/rest/refund.do"
	GetOrderStatusExtended string = "/payment/rest/getOrderStatusExtended.do"
	GetReceiptStatus       string = "/payment/rest/getReceiptStatus.do"
	UnBindCard             string = "/payment/rest/unBindCard.do"
	BindCard               string = "/payment/rest/bindCard.do"
	GetBindings            string = "/payment/rest/getBindings.do"
	ExtendBinding          string = "/payment/rest/extendBinding.do"
	ApplePay               string = "/payment/applepay/payment.do"
	SamsungPay             string = "/payment/samsung/payment.do"
	GooglePay              string = "/payment/google/payment.do"
	VerifyEnrollment       string = "/payment/rest/verifyEnrollment.do"
	UpdateSSLCardList      string = "/payment/rest/updateSSLCardList.do"
)
