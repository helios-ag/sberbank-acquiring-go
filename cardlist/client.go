package cardlist

import (
	"context"
	"encoding/json"
	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
	"net/http"
)

type Client struct {
	API acquiring.API
}

// UpdateSSLCardList request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:updateSSLCardList
func (c *Client) UpdateSSLCardList(ctx context.Context, mdorder string, jsonParams map[string]string) (*schema.Response, *http.Response, error) {
	path := endpoints.UpdateSSLCardList
	body := make(map[string]string)
	body["mdorder"] = mdorder

	var response schema.Response
	req, err := c.API.NewRestRequest(ctx, "GET", path, body, jsonParams)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &response)
	if err != nil {
		return nil, result, err
	}
	_ = json.NewDecoder(result.Body).Decode(&response)

	return &response, result, err
}
