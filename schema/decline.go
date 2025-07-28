package schema

// DeclineResponse is mapped Decline response from Decline request
type DeclineResponse struct {
	ErrorCode    int    `json:"errorCode,string"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
