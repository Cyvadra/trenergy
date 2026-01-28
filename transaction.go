package trenergy

import (
	"context"
	"fmt"
	"net/url"
)

// InternalTransaction represents an internal transaction.
type InternalTransaction struct {
	ID             int         `json:"id"`
	Amount         float64     `json:"amount"`
	Type           int         `json:"type"`
	Coin           int         `json:"coin"`            // 1 - main, 3 - energy?
	InstantBalance interface{} `json:"instant_balance"` // 0 or null?
	TxableID       int         `json:"txable_id"`
	CreatedAt      string      `json:"created_at"`
}

// InternalTransactionParams for filtering.
type InternalTransactionParams struct {
	Page          int
	PerPage       int
	Types         []int
	Date          string // Y-m-d
	SortBy        string // created_at/amount
	SortDirection string // asc/desc
}

// GetInternalTransactions retrieves internal transactions.
func (c *Client) GetInternalTransactions(ctx context.Context, params InternalTransactionParams) (*APIResponse[[]InternalTransaction], error) {
	u, _ := url.Parse("/api/transactions/internal")
	q := u.Query()

	if params.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PerPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", params.PerPage))
	}
	for _, t := range params.Types {
		q.Add("types[]", fmt.Sprintf("%d", t))
	}
	if params.Date != "" {
		q.Set("date", params.Date)
	}
	if params.SortBy != "" {
		q.Set("sort_by", params.SortBy)
	}
	if params.SortDirection != "" {
		q.Set("sort_direction", params.SortDirection)
	}
	u.RawQuery = q.Encode()

	var resp APIResponse[[]InternalTransaction]
	err := c.sendRequest(ctx, "GET", u.String(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
