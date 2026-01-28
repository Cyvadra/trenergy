# Tr.Energy Go SDK

A Go client library for the [Tr.Energy](https://tr.energy) API, enabling easy integration with the TRON Energy marketplace.

## Installation

```bash
go get github.com/cyvadra/trenergy
```

## Usage

### Initialization

Initialize the client with your API key. You can find your API key in your Tr.Energy account settings.

```go
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

	// ...
}
```

### Testnet Support (Nile)

To use the Nile Testnet, use the `WithTestNet()` option. This will automatically use the testnet base URL and a default test key if none is provided.

```go
client := trenergy.NewClient("", trenergy.WithTestNet())
```

### Get Account Info

Retrieve your account information, including balance and status.

```go
account, err := client.GetAccountInfo(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User: %s, Balance: %f\n", account.Data.Name, account.Data.Balance)
```

### Activate Address

Activate a TRON address if it's inactive (not on-chain yet).

```go
resp, err := client.ActivateAddress(context.Background(), "TCsKjXeT652tbDhjVzVdFEtacXNwhBCcDQ")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Activation Status: %v\n", resp.Status)
```

### Create Energy Order (Bootstrap)

Purchase energy for an address using a "bootstrap" order.

```go
order, err := client.CreateBootstrapOrder(context.Background(), trenergy.ConsumerParams{
    Address:        "TCsKjXeT652tbDhjVzVdFEtacXNwhBCcDQ",
    PaymentPeriod:  15, // e.g., 1 hour (check API docs for allowed periods)
    AutoRenewal:    false,
    ResourceAmount: 65150, // Amount of energy
    Resource:       1,     // 1 for ENERGY, 0 for BANDWIDTH
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order Created: ID %d\n", order.Data.ID)
```

### List Consumers

List your current energy consumers.

```go
consumers, err := client.ListConsumers(context.Background(), 1) // 1 is the page number
if err != nil {
    log.Fatal(err)
}
for _, c := range consumers.Data {
    fmt.Printf("Consumer: %s (%s)\n", c.Name, c.Address)
}
```

## Features

- **Account Management**: Check balance and account details.
- **Consumer Management**: Create, list, delete, and manage energy consumers.
- **Order Management**: Create bootstrap orders for immediate energy.
- **Wallet & Transactions**: (Helper functions for Tron wallet interactions if available in SDK).
- **Testnet Support**: Seamlessly switch to Nile Testnet for development.

## License

MIT
