package trenergy

import (
	"context"
	"fmt"
)

// Wallet represents a user's wallet.
type Wallet struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListWallets retrieves a list of wallets.
func (c *Client) ListWallets(ctx context.Context) (*APIResponse[[]Wallet], error) {
	var resp APIResponse[[]Wallet]
	err := c.sendRequest(ctx, "GET", "/api/wallets", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddWallet adds a new wallet.
func (c *Client) AddWallet(ctx context.Context, address string) (*APIResponse[*Wallet], error) {
	data := make(map[string]string)
	data["address"] = address

	// Sample 2691 uses FormData
	var resp APIResponse[*Wallet]
	err := c.postMultipart(ctx, "/api/wallets", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteWallet deletes a wallet.
func (c *Client) DeleteWallet(ctx context.Context, id int) (*APIResponse[struct{}], error) {
	path := fmt.Sprintf("/api/wallets/%d", id)
	var resp APIResponse[struct{}]
	err := c.sendRequest(ctx, "DELETE", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
