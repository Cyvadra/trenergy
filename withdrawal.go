package trenergy

import (
	"context"
	"fmt"
)

// Withdrawal represents a withdrawal record.
type Withdrawal struct {
	ID        int     `json:"id"`
	TrxAmount float64 `json:"trx_amount"`
	Status    string  `json:"status"`
	Address   string  `json:"address"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// ListWithdrawals retrieves a list of withdrawals.
func (c *Client) ListWithdrawals(ctx context.Context, page int) (*APIResponse[[]Withdrawal], error) {
	path := fmt.Sprintf("/api/withdrawals?page=%d", page)
	var resp APIResponse[[]Withdrawal]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateWithdrawal creates a new withdrawal request.
func (c *Client) CreateWithdrawal(ctx context.Context, amount float64, address string, otp string) (*APIResponse[struct{}], error) {
	data := make(map[string]string)
	data["trx_amount"] = fmt.Sprintf("%f", amount)
	if address != "" {
		data["address"] = address
	}
	if otp != "" {
		data["one_time_password"] = otp
	}

	var resp APIResponse[struct{}]
	err := c.postMultipart(ctx, "/api/withdrawals", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
