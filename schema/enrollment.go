package schema

// EnrollmentResponse is response mapping from Enrollment request
type EnrollmentResponse struct {
	ErrorCode          int    `json:"errorCode,omitempty"`
	ErrorMessage       string `json:"errorMessage,omitempty"`
	Enrolled           byte   `json:"enrollment,omitempty"`
	EmitterName        string `json:"emitterName,omitempty"`
	EmitterCountryCode string `json:"emitterCountryCode,omitempty"`
}
