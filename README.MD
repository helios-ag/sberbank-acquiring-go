[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Go](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml/badge.svg)](https://github.com/helios-ag/sberbank-acquiring-go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/helios-ag/sberbank-acquiring-go/branch/master/graph/badge.svg)](https://codecov.io/gh/helios-ag/sberbank-acquiring-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/helios-ag/sberbank-acquiring-go)](https://goreportcard.com/report/github.com/helios-ag/sberbank-acquiring-go)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/helios-ag/sberbank-acquiring-go)
# Sberbank Acquiring API Wrapper

## Описание

## Пример

```golang
package main

import (
	"context"
	"fmt"
	"github.com/helios-ag/sberbank-acquiring-go/acquiring"
	"github.com/helios-ag/sberbank-acquiring-go/currency"
)

func main() {
	cfg := acquiring.ClientConfig{
		UserName:           "test-api", // Replace with your own
		Currency:           currency.RUB,
		Password:           "test", // Replace with your own
		Language:           "ru",
		SessionTimeoutSecs: 1200,
		SandboxMode:        true,
	}

	client, err := acquiring.NewClient(&cfg)
	if err != nil {
		panic(err)
	}
	order := acquiring.Order{
		OrderNumber: "test",
		Amount:      100,
		Description: "My Order for Client",
	}
	result, _, err := client.RegisterOrder(context.Background(), order)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.ErrorCode)
	fmt.Println(result.ErrorMessage)
	fmt.Println(result.FormUrl)
	fmt.Println(result.OrderId)
}

```

