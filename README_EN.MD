[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Go](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml/badge.svg)](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/helios-ag/sberbank-acquiring-go/graph/badge.svg?token=eCWqf8zkm3)](https://codecov.io/gh/helios-ag/sberbank-acquiring-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/helios-ag/sberbank-acquiring-go)](https://goreportcard.com/report/github.com/helios-ag/sberbank-acquiring-go)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/helios-ag/sberbank-acquiring-go)
# Sberbank Acquiring API Wrapper

[Sberbank Acquiring API](https://securepayments.sberbank.ru/wiki/doku.php/integration:api:start) Библиотека предоставляет удобную оболочку над REST‑интерфейсом Sberbank Acquiring.

Она позволяет быстро зарегистрировать заказ, выполнить платеж, проверить статус, оформить возврат, работать с привязками карт и использовать мобильные способы оплаты (Apple Pay, Google Pay, Mir Pay и т.д.).

**Возможности**

* Регистрация заказов: register и registerPreAuth.

* Управление платежами: депонирование (deposit), сторнирование (reverse), возвраты (refund и instantRefund).

* Запрос статуса заказа и чека.

* Работа с привязками карт: создание, продление, удаление и получение списка.

* Проверка карты на участие в 3‑D Secure (verifyEnrollment).

* Получение внешнего чека.

* Мобильные платежи: Apple Pay, Google Pay, Samsung Pay, Mir Pay, Mir Pay Direct.

* Обновление SSL‑списка карт.

* Утилиты и схемы ответов для типовых операций.

## Установка


Используйте Go modules (go mod init)

``` sh
go mod init
go get -u github.com/helios-ag/sberbank-acquiring-go
```

Импортируем зависимости

``` go
import (
	"github.com/helios-ag/sberbank-acquiring-go/acquiring"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
)
```

```golang
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
