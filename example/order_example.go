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
		UserName:           "test-api", // Replace with your own
		Currency:           currency.RUB,
		Password:           "test", // Replace with your own
		Language:           "ru",
		SessionTimeoutSecs: 1200,
		SandboxMode:        true,
	}

	acquiring.SetConfig(cfg)

	order := orders.Order{
		OrderNumber: "test",
		Amount:      100,
		Description: "My Order for Client",
	}
	result, _, err := orders.RegisterOrder(context.Background(), order)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.ErrorCode)
	fmt.Println(result.ErrorMessage)
	fmt.Println(result.FormUrl)
	fmt.Println(result.OrderId)
}
