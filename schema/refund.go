package schema

// InstantRefundResponse — структура ответа от метода instantRefund.do
type InstantRefundResponse struct {
	ErrorCode    int    `json:"errorCode,string"`    // Код ошибки
	ErrorMessage string `json:"errorMessage,string"` // Описание ошибки
	OrderID      string `json:"orderId"`             // Номер заказа в платёжной системе
	OrderStatus  int    `json:"orderStatus"`         // Статус заказа

	// Дополнительно, если возвращается блок orderStatus
	ApprovalCode string `json:"approvalCode"` // Код авторизации МПС
	RRN          string `json:"rrn"`          // Reference Retrieval Number
}

// ProcessRawSumRefundResponse — структура ответа от метода processRawSumRefund.do
type ProcessRawSumRefundResponse struct {
	ErrorCode    int    `json:"errorCode,string"`    // Код ошибки
	ErrorMessage string `json:"errorMessage,string"` // Описание ошибки
}
