package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cyvadra/trenergy"
)

func main() {
	apiKey := os.Getenv("TRENERGY_API_KEY")
	client := trenergy.NewClient(apiKey)
	// Get Account Info
	account, err := client.GetAccountInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %s, Balance: %f\n", account.Data.Name, account.Data.Balance)
	// Activate Address
	bootstrapOrder, err := client.CreateBootstrapOrder(context.Background(), trenergy.ConsumerParams{
		Address:        "TCsKjXeT652tbDhjVzVdFEtacXNwhBCcDQ",
		PaymentPeriod:  15,
		AutoRenewal:    false,
		ResourceAmount: 65150,
		Resource:       1, // 0 - BANDWIDTH, 1 - ENERGY
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bootstrap Order: %v\n", bootstrapOrder)
	// List Consumers
	consumers, err := client.ListConsumers(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range consumers.Data {
		fmt.Printf("Consumer: %s (%s)\n", c.Name, c.Address)
	}
}
