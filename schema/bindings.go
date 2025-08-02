package schema

// BindingsResponse is mapped Binding response from GetBindings request
type BindingsResponse struct {
	ErrorCode    int    `json:"errorCode,string"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Bindings     []struct {
		BindingId  string `json:"bindingId,omitempty"`
		MaskedPan  string `json:"maskedPan,omitempty"`
		ExpiryDate string `json:"expiryDate,omitempty"`
	} `json:"bindings,omitempty"`
}

// BindingsByCardOrIdResponse is mapped Binding response from GetBindingsByCardOrId request
type BindingsByCardOrIdResponse struct {
	ErrorCode    int    `json:"errorCode,string"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Bindings     []struct {
		BindingId  string `json:"bindingId,omitempty"`
		MaskedPan  string `json:"maskedPan,omitempty"`
		ExpiryDate string `json:"expiryDate,omitempty"`
		ClientId   string `json:"clientId,omitempty"`
	} `json:"bindings,omitempty"`
}
type BindingsNoPaymentResponse struct {
	ErrorCode      int    `json:"errorCode,string"`
	Error          string `json:"error,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	BindingId      string `json:"bindingId,omitempty"`
	CardholderName string `json:"cardholderName,omitempty"`
	ClientId       string `json:"clientId,omitempty"`
	ExpiryDate     string `json:"expiryDate,omitempty"`
	MaskedPan      string `json:"maskedPan,omitempty"`
}
