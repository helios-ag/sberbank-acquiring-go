[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Go](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml/badge.svg)](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/helios-ag/sberbank-acquiring-go/graph/badge.svg?token=eCWqf8zkm3)](https://codecov.io/gh/helios-ag/sberbank-acquiring-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/helios-ag/sberbank-acquiring-go)](https://goreportcard.com/report/github.com/helios-ag/sberbank-acquiring-go)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/helios-ag/sberbank-acquiring-go)
# sberbank‑acquiring‑go

Библиотека предоставляет удобную оболочку над [REST‑интерфейсом Sberbank Acquiring](https://securepayments.sberbank.ru/wiki/doku.php/integration:api:start).  
Она позволяет быстро зарегистрировать заказ, выполнить платеж, проверить статус, оформить возврат, работать с привязками карт и использовать мобильные способы оплаты (Apple Pay, Google Pay, Mir Pay и т.д.).

## Возможности

- Регистрация заказов: `register` и `registerPreAuth`.
- Управление платежами: Запрос завершения оплаты заказа (`deposit`), Запрос отмены оплаты заказа (`reverse`), Запрос возврата средств оплаты заказа (`refund` и `instantRefund`).
- Запрос статуса заказа и чека.
- Работа с привязками карт.
- Проверка карты на участие в 3‑D Secure (`verifyEnrollment`).
- Мобильные платежи: Apple Pay, Google Pay, Samsung Pay, Mir Pay, Mir Pay Direct.
- Утилиты и схемы ответов для типовых операций.

## Установка

```bash
go get -u github.com/helios-ag/sberbank-acquiring-go
```

Используйте **Go modules** (`go mod init`), чтобы зависимость подключалась автоматически.

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"

    acquiring "github.com/helios-ag/sberbank-acquiring-go"
    "github.com/helios-ag/sberbank-acquiring-go/currency"
    "github.com/helios-ag/sberbank-acquiring-go/orders"
)

func main() {
    cfg := acquiring.ClientConfig{
        UserName:           "test-api",       // ваши учётные данные
        Password:           "test",
        Currency:           currency.RUB,
        Language:           "ru",
        SessionTimeoutSecs: 1200,
        SandboxMode:        true,             // sandbox: true, production: false
    }
    acquiring.SetConfig(cfg)

    order := orders.Order{
        OrderNumber: "order-001",
        Amount:      10000,                   // сумма в копейках
        Description: "Покупка товара",
    }

    resp, _, err := orders.RegisterOrder(context.Background(), order)
    if err != nil {
        panic(err)
    }

    fmt.Println("OrderID:", resp.OrderId)
    fmt.Println("Payment form:", resp.FormUrl)
}
```

## Работа с заказами

### Получение статуса заказа

```go
statusResp, _, err := orders.GetOrderStatus(context.Background(),
    orders.Order{OrderNumber: "order-001"})
if err != nil {
    panic(err)
}
fmt.Println("Status:", statusResp.OrderStatus)
```

### Возврат средств

```go
refundResp, _, err := orders.RefundOrder(context.Background(),
    orders.Order{OrderNumber: "order-001", Amount: 5000})
if err != nil {
    panic(err)
}
fmt.Println("Refund error:", refundResp.ErrorMessage)
```

## Привязка карт

```go
import "github.com/helios-ag/sberbank-acquiring-go/bind"

bindResp, _, err := bind.BindCard(context.Background(),
    bind.Binding{BindingID: "id123"})
if err != nil {
    panic(err)
}
fmt.Println("Result:", bindResp.ErrorCode)
```

Получение списка привязок:

```go
bindings, _, err := bind.GetBindings(context.Background(), "client123", nil)
if err != nil {
    panic(err)
}
fmt.Println("Bindings:", bindings.Binding)
```

## Запрос проверки вовлечённости карты в 3DS

```go
import "github.com/helios-ag/sberbank-acquiring-go/enrollment"

enrollResp, _, err := enrollment.VerifyEnrollment(context.Background(), "4111111111111111")
if err != nil {
    panic(err)
}
fmt.Println("Enrolled:", enrollResp.Enrolled)
```

## Мобильные платежи

### Apple Pay

```go
import "github.com/helios-ag/sberbank-acquiring-go/mobile"

appleReq := mobile.ApplePaymentRequest{
    OrderNumber:  "order-002",
    Merchant:     "merchant.test",
    PaymentToken: "base64Token",
}
appleResp, _, err := mobile.PayWithApplePay(context.Background(), appleReq)
if err != nil {
    panic(err)
}
fmt.Println("ApplePay status:", appleResp.ErrorCode)
```

### Google Pay

```go
googleReq := mobile.GooglePaymentRequest{
    OrderNumber:  "order-003",
    Merchant:     "merchant.test",
    PaymentToken: "base64Token",
    Language:     "ru",
    ReturnUrl:    "https://example.org/return",
}
googleResp, _, err := mobile.PayWithGooglePay(context.Background(), googleReq)
if err != nil {
    panic(err)
}
fmt.Println("GooglePay status:", googleResp.ErrorCode)
```

Аналогично доступны методы для Mir Pay и Mir Pay Direct.

## Запрос возврата средств по устаревшему заказу сроком более двух лет

```go
import "github.com/helios-ag/sberbank-acquiring-go/instant_refund"

req := instant_refund.InstantRefundRequest{
    UserName:    "user",
    Password:    "pass",
    Amount:      10000,
    OrderNumber: ptr.String("order-004"),
}
resp, _, err := instant_refund.InstantRefund(context.Background(), req)
if err != nil {
    panic(err)
}
fmt.Println("Instant refund:", resp.ErrorCode)
```

## Запрос передачи данных чека

```go
import "github.com/helios-ag/sberbank-acquiring-go/external_receipt"

receiptReq := external_receipt.ExternalReceiptRequest{
    UserName: "user",
    Password: "pass",
    MdOrder:  "order-001",
    Receipt:  &external_receipt.Receipt{PaymentType: 1},
}
receipt, _, err := external_receipt.GetExternalReceipt(context.Background(), receiptReq)
if err != nil {
    panic(err)
}
fmt.Println("Receipt status:", receipt.Status)
```

## Запрос добавления карты в список SSL-карт

```go
import "github.com/helios-ag/sberbank-acquiring-go/cardlist"

resp, _, err := cardlist.UpdateSSLCardList(context.Background(), "mdOrder123", nil)
if err != nil {
    panic(err)
}
fmt.Println("Update card list:", resp.ErrorCode)
```

## Лицензия

Проект распространяется по лицензии [MIT](LICENSE).