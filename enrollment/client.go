package enrollment

import (
	"context"
	"encoding/json"
	"fmt"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"net/http"
)

type Client struct {
	API acquiring.API
}

// VerifyEnrollment request
// Checks if card enrolled in 3D Sec
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:verifyEnrollment
func VerifyEnrollment(ctx context.Context, pan string) (*schema.EnrollmentResponse, *http.Response, error) {
	return getClient().VerifyEnrollment(ctx, pan)
}

// VerifyEnrollment request
// Checks if card enrolled in 3D Sec
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:verifyEnrollment
func (c Client) VerifyEnrollment(ctx context.Context, pan string) (*schema.EnrollmentResponse, *http.Response, error) {
	path := endpoints.VerifyEnrollment

	if err := validatePan(pan); err != nil {
		return nil, nil, err
	}

	body := make(map[string]string)
	body["pan"] = pan

	var enrollmentResponse schema.EnrollmentResponse
	req, err := c.API.NewRestRequest(ctx, http.MethodGet, path, body, nil)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &enrollmentResponse)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&enrollmentResponse)

	return &enrollmentResponse, result, err
}

func validatePan(pan string) error {
	panLength := len(pan)
	if panLength < 13 || panLength > 19 {
		return fmt.Errorf("pan length shouldn't be less 13 or more 19 symbols")
	}
	return nil
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
