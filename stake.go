package trenergy

import (
	"context"
	"fmt"
)

// Stake represents a stake record.
type Stake struct {
	ID          int     `json:"id"`
	Resource    string  `json:"resource"`
	TrxAmount   float64 `json:"trx_amount"` // Sample has it as number.
	Type        int     `json:"type"`
	IsCloses    bool    `json:"is_closes"`
	ClosesAt    *string `json:"closes_at"`
	DefrostedAt *string `json:"defrosted_at"`
	CreatedAt   string  `json:"created_at"`
	AvailableAt string  `json:"available_at"`
}

// ListStakes retrieves a list of stakes.
func (c *Client) ListStakes(ctx context.Context, page int) (*APIResponse[[]Stake], error) {
	path := fmt.Sprintf("/api/stakes?page=%d", page)
	var resp APIResponse[[]Stake]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateStakeParams
type CreateStakeParams struct {
	TrxAmount float64
}

// CreateStake creates a new stake.
func (c *Client) CreateStake(ctx context.Context, params CreateStakeParams) (*APIResponse[struct{}], error) {
	data := make(map[string]string)
	data["trx_amount"] = fmt.Sprintf("%f", params.TrxAmount)

	var resp APIResponse[struct{}]
	err := c.postMultipart(ctx, "/api/stakes", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnstakeParams
type UnstakeParams struct {
	TrxAmount float64
	OTP       string
}

// Unstake unstakes an amount.
func (c *Client) Unstake(ctx context.Context, params UnstakeParams) (*APIResponse[struct {
	UnstakeDate string `json:"unstake_date"`
}], error) {
	data := make(map[string]string)
	data["trx_amount"] = fmt.Sprintf("%f", params.TrxAmount)
	if params.OTP != "" {
		data["one_time_password"] = params.OTP
	}

	var resp APIResponse[struct {
		UnstakeDate string `json:"unstake_date"`
	}]
	err := c.postMultipart(ctx, "/api/stakes/unstake", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SyncStakes syncs stakes from the blockchain.
func (c *Client) SyncStakes(ctx context.Context) (*APIResponse[struct{}], error) {
	var resp APIResponse[struct{}]
	err := c.sendRequest(ctx, "POST", "/api/stakes/sync", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// StakeProfitabilityItem represents profitability data point.
type StakeProfitabilityItem struct {
	Received float64 `json:"received"`
	Date     string  `json:"date"`
}

// GetStakeProfitability retrieves stake profitability info.
func (c *Client) GetStakeProfitability(ctx context.Context, period int) (*APIResponse[[]StakeProfitabilityItem], error) {
	// period can be 7, 30, 365
	path := fmt.Sprintf("/api/stakes/profitability?period=%d", period)
	var resp APIResponse[[]StakeProfitabilityItem]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
