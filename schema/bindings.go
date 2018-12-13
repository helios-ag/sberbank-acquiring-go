package schema

// BindingsResponse is mapped Binding response from GetBindings request
type BindingsResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Bindings     []struct {
		BindingId  string `json:"bindingId,omitempty"`
		MaskedPan  string `json:"maskedPan,omitempty"`
		ExpiryDate string `json:"expiryDate,omitempty"`
	} `json:"bindings,omitempty"`
}
