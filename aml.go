package trenergy

import (
	"context"
	"fmt"
)

// AMLCheck represents an AML check result.
type AMLCheck struct {
	Address    string           `json:"address"`
	Blockchain *string          `json:"blockchain"`
	TxID       *string          `json:"txid"`
	Status     string           `json:"status"` // e.g. "completed"
	Context    *AMLCheckContext `json:"context"`
	CreatedAt  string           `json:"created_at"`
}

type AMLCheckContext struct {
	Pending   bool        `json:"pending"`
	Entities  []AMLEntity `json:"entities"`
	RiskScore float64     `json:"riskScore"`
}

type AMLEntity struct {
	Level     string  `json:"level"`
	Entity    string  `json:"entity"`
	RiskScore float64 `json:"riskScore"`
}

// ListAMLChecks retrieves a list of AML checks.
func (c *Client) ListAMLChecks(ctx context.Context, page int) (*APIResponse[[]AMLCheck], error) {
	path := fmt.Sprintf("/api/aml?page=%d", page)
	var resp APIResponse[[]AMLCheck]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckAML performs a new AML check.
func (c *Client) CheckAML(ctx context.Context, address string, txid string) (*APIResponse[*AMLCheck], error) {
	data := make(map[string]string)
	data["address"] = address
	if txid != "" {
		data["txid"] = txid
	}

	var resp APIResponse[*AMLCheck]
	err := c.postMultipart(ctx, "/api/aml/check", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAMLCheck retrieves a specific AML check.
func (c *Client) GetAMLCheck(ctx context.Context, id int) (*APIResponse[*AMLCheck], error) {
	path := fmt.Sprintf("/api/aml/%d", id)
	var resp APIResponse[*AMLCheck]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
