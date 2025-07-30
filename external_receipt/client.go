package receipt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	acquiring "github.com/helios-ag/sberbank-acquiring-go"
	"github.com/helios-ag/sberbank-acquiring-go/endpoints"
	"github.com/helios-ag/sberbank-acquiring-go/schema"
)

type Client struct {
	API acquiring.API
}

// ExternalReceiptRequest is used for building GetExternalReceipt request
//
// Language - язык в кодировке ISO 639-1. Если не указан — будет использован язык по умолчанию.
// UserName - логин служебной учётной записи продавца.
// Password - пароль служебной учётной записи продавца.
// MdOrder - уникальный номер заказа в платёжном шлюзе.
// Receipt - блок с параметрами чека.
// CashboxID - идентификатор кассы.
// BasketID - идентификатор корзины покупки или возврата.
// JsonParams - Дополнительные параметры запроса.
type ExternalReceiptRequest struct {
	Language   *string     `json:"language,omitempty"`   // Язык в кодировке ISO 639-1
	UserName   string      `json:"userName"`             // Логин служебной учётной записи продавца
	Password   string      `json:"password"`             // Пароль служебной учётной записи продавца
	MdOrder    string      `json:"mdOrder"`              // Уникальный номер заказа в платёжном шлюзе
	Receipt    *Receipt    `json:"receipt"`              // Блок с параметрами чека (структура Receipt уже определена)
	JSONParams *JSONParams `json:"jsonParams,omitempty"` // Дополнительные параметры запроса
}

type JSONParams struct {
	CashboxID *int64 `json:"cashboxId,omitempty"` // Идентификатор кассы
	BasketID  *int64 `json:"basketId,omitempty"`  // Идентификатор корзины покупки или возврата
}

type Receipt struct {
	Type                    *int     `json:"type,omitempty"`                      // Признак способа оплаты: 0 - приход, 1 - возврат
	PaymentType             int      `json:"paymentType"`                         // Способ оплаты: 1 - оплата, 3 - возврат
	FnNumber                *string  `json:"fn_number,omitempty"`                 // Номер фискального накопителя
	FiscalDocumentNumber    *int64   `json:"fiscal_document_number,omitempty"`    // Фискальный номер документа
	FiscalDocumentAttribute *string  `json:"fiscal_document_attribute,omitempty"` // Фискальный признак документа
	AmountTotal             *float64 `json:"amount_total,omitempty"`              // Итоговая сумма чека (до 15 целых и 2 знаков после запятой)
	ReceiptDateTime         *string  `json:"receipt_date_time,omitempty"`         // Дата и время чека (формат: yyyy:MM:dd HH:mm:ss)
}

// GetExternalReceipt request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:externalreceipt
func GetExternalReceipt(ctx context.Context, externalReceipt ExternalReceiptRequest) (*schema.ExternalReceipt, *http.Response, error) {
	return getClient().GetExternalReceipt(ctx, externalReceipt)
}

// GetExternalReceipt GetExternalReceipt request
// see https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:requests:externalreceipt
func (c Client) GetExternalReceipt(ctx context.Context, externalReceipt ExternalReceiptRequest) (*schema.ExternalReceipt, *http.Response, error) {
	path := endpoints.ExternalReceipt

	if err := validateExternalReceiptRequest(externalReceipt); err != nil {
		return nil, nil, err
	}

	var receipt, _ = json.Marshal(externalReceipt.Receipt)
	var jsonParams, _ = json.Marshal(externalReceipt.JSONParams)
	body := map[string]string{
		"userName": externalReceipt.UserName,
		"password": externalReceipt.Password,
		"mdOrder":  externalReceipt.MdOrder,
		"receipt":  string(receipt[:]),
	}
	if externalReceipt.Language != nil {
		body["language"] = *externalReceipt.Language
	}

	if externalReceipt.JSONParams != nil {
		body["jsonParams"] = string(jsonParams[:])
	}
	var response schema.ExternalReceipt

	req, err := c.API.NewRequest(ctx, http.MethodPost, path, body)

	if err != nil {
		return nil, nil, err
	}
	result, err := c.API.Do(req, &response)
	if err != nil {
		return nil, result, err
	}

	return &response, result, err
}

func validateExternalReceiptRequest(externalReceiptRequest ExternalReceiptRequest) error {
	if externalReceiptRequest.UserName == "" || externalReceiptRequest.Password == "" || externalReceiptRequest.MdOrder == "" || externalReceiptRequest.Receipt == nil {
		return fmt.Errorf("userName and Password and mdOrder and Receipt are required")
	}

	return nil
}

func getClient() Client {
	return Client{acquiring.GetAPI()}
}
