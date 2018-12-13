package schema

// ReceiptStatus is response received from GetReceiptStatus
type ReceiptStatus struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	OrderNumber  string `json:"orderNumber,omitempty"`
	OrderId      string `json:"orderId,omitempty"`
	DaemonCode   string `json:"daemonCode,omitempty"`
	DeviceCode   string `json:"deviceCode,omitempty"`
	Receipt      []struct {
		ReceiptStatus           int    `json:"receiptStatus,omitempty"`
		Uuid                    string `json:"uuid,omitempty"`
		ShiftNumber             int    `json:"shift_number,omitempty"`
		ReceiptNumber           int    `json:"receipt_number,omitempty"`
		ReceiptDatetime         string `json:"receipt_datetime,omitempty"`
		FnNumber                string `json:"fn_number,omitempty"`
		DeviceNumber            string `json:"device_number,omitempty"`
		FiscalDocumentNumber    int    `json:"fiscal_document_number,omitempty"`
		FiscalDocumentAttribute string `json:"fiscal_document_attribute,omitempty"`
		AmountTotal             string `json:"amount_total,omitempty"`
		SerialNumber            string `json:"serial_number,omitempty"`
		FnsSite                 string `json:"fnsSite,omitempty"`
		OFD                     struct {
			Name    string `json:"name,omitempty"`
			Website string `json:"website,omitempty"`
			INN     string `json:"receipt_number,omitempty"`
		} `json:"OFD,omitempty"`
	} `json:"receipt"`
}
