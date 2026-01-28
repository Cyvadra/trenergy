package trenergy

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Consumer represents a consumer entity.
type Consumer struct {
	ID                    int         `json:"id"`
	Name                  string      `json:"name"`
	Address               string      `json:"address"`
	Resource              string      `json:"resource"`
	ResourceAmount        interface{} `json:"resource_amount"`         // Can be string or int based on samples
	DesiredResourceAmount interface{} `json:"desired_resource_amount"` // Can be string
	CreationType          int         `json:"creation_type"`
	PaymentPeriod         int         `json:"payment_period"`
	AutoRenewal           bool        `json:"auto_renewal"`
	IsActive              bool        `json:"is_active"`
	Order                 *Order      `json:"order"`
	WebhookURL            string      `json:"webhook_url"`
	CreatedAt             string      `json:"created_at"`
	UpdatedAt             string      `json:"updated_at"`
	EstimatedCostTrx      float64     `json:"estimated_cost_trx,omitempty"`
}

type Order struct {
	Status               int    `json:"status"`
	CompletionPercentage int    `json:"completion_percentage"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	ValidUntil           string `json:"valid_until"`
}

// ConsumerParams are parameters for creating/updating a consumer.
type ConsumerParams struct {
	Address        string
	PaymentPeriod  int // usage: 15, 60, etc.
	AutoRenewal    bool
	ResourceAmount int64
	Name           string
	Resource       int // 0 - BANDWIDTH, 1 - ENERGY
	WebhookURL     string
}

// CreateBootstrapOrder creates a new consumer order.
func (c *Client) CreateBootstrapOrder(ctx context.Context, params ConsumerParams) (*APIResponse[*Consumer], error) {
	// Note: samples show form-data
	data := make(map[string]string)
	data["address"] = params.Address
	data["payment_period"] = strconv.Itoa(params.PaymentPeriod)
	if params.AutoRenewal {
		data["auto_renewal"] = "1"
	} else {
		data["auto_renewal"] = "0"
	}
	data["resource_amount"] = strconv.FormatInt(params.ResourceAmount, 10)

	if params.Name != "" {
		data["name"] = params.Name
	}
	// Resource default is Energy (1), optional? Sample sends "1".
	data["resource"] = strconv.Itoa(params.Resource)

	if params.WebhookURL != "" {
		data["webhook_url"] = params.WebhookURL
	}

	var resp APIResponse[*Consumer]
	err := c.postMultipart(ctx, "/api/consumers/bootstrap-order", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListConsumers returns a list of consumers with pagination.
// Using explicit params for pagination for simplicity.
func (c *Client) ListConsumers(ctx context.Context, page int) (*APIResponse[[]Consumer], error) {
	path := fmt.Sprintf("/api/consumers?page=%d", page)
	var resp APIResponse[[]Consumer]
	// Note: samples show GET
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConsumer retrieves a single consumer by ID.
func (c *Client) GetConsumer(ctx context.Context, id int) (*APIResponse[*Consumer], error) {
	path := fmt.Sprintf("/api/consumers/%d", id)
	var resp APIResponse[*Consumer]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ActivateConsumer activates a consumer.
func (c *Client) ActivateConsumer(ctx context.Context, id int) (*APIResponse[struct{}], error) {
	path := fmt.Sprintf("/api/consumers/%d/activate", id)
	var resp APIResponse[struct{}]
	err := c.sendRequest(ctx, "POST", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeactivateConsumer deactivates a consumer.
func (c *Client) DeactivateConsumer(ctx context.Context, id int) (*APIResponse[struct{}], error) {
	path := fmt.Sprintf("/api/consumers/%d/deactivate", id)
	var resp APIResponse[struct{}]
	err := c.sendRequest(ctx, "POST", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateConsumer updates a consumer.
func (c *Client) UpdateConsumer(ctx context.Context, id int, params ConsumerParams) (*APIResponse[struct{}], error) {
	path := fmt.Sprintf("/api/consumers/%d", id)
	// Sample uses PATCH and form-urlencoded
	data := url.Values{}
	if params.Name != "" {
		data.Set("name", params.Name)
	}
	// Add other fields if update allows them. Sample only shows name being updated, but docs might allow more.
	// Assuming logic similar to create for other optional fields if documentation implies.
	// For now implementing what is in samples (only name is shown in PATCH sample).

	// BUT sample 801 also shows "resource_amount", "payment_period" disabled... implying they might be updateable or just visible in docs.
	// Let's assume we can update what we can.

	var resp APIResponse[struct{}]
	// We need a helper for PATCH with urlencoded.
	// Since postForm is POST, we need something generic or just construct it here.

	// Let's reuse postForm logic but with PATCH
	u, _ := url.Parse(path)
	finalURL := c.baseURL.ResolveReference(u)

	req, err := http.NewRequestWithContext(ctx, "PATCH", finalURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range c.headers {
		req.Header[k] = v
	}

	_, err = c.Do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteConsumer deletes a consumer.
func (c *Client) DeleteConsumer(ctx context.Context, id int) (*APIResponse[struct{}], error) {
	path := fmt.Sprintf("/api/consumers/%d", id)
	var resp APIResponse[struct{}]
	err := c.sendRequest(ctx, "DELETE", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ConsumerPayment represents a payment history item.
type ConsumerPayment struct {
	Amount    float64     `json:"amount"`
	Type      int         `json:"type"`
	Coin      int         `json:"coin"`
	Quantity  interface{} `json:"quantity"` // Can be int or null
	CreatedAt string      `json:"created_at"`
}

// GetConsumerPayments gets payments for a consumer.
func (c *Client) GetConsumerPayments(ctx context.Context, id int, page int) (*APIResponse[[]ConsumerPayment], error) {
	path := fmt.Sprintf("/api/consumers/%d/payments?page=%d", id, page)
	var resp APIResponse[[]ConsumerPayment]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ActivateAddressData represents the data for activate address response.
type ActivateAddressData struct{}

// ActivateAddress activates a generic address.
func (c *Client) ActivateAddress(ctx context.Context, address string) (*APIResponse[ActivateAddressData], error) {
	// Sample 1150 POST /api/extra/activate-address with formdata
	data := make(map[string]string)
	data["address"] = address

	var resp APIResponse[ActivateAddressData]
	err := c.postMultipart(ctx, "/api/extra/activate-address", data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ConsumersSummary represents summary stats.
type ConsumersSummary struct {
	TotalCount                 int                `json:"total_count"`
	TotalEnergyCount           int                `json:"total_energy_count"`
	TotalBandwidthCount        int                `json:"total_bandwidth_count"`
	ActiveCount                int                `json:"active_count"`
	ActiveEnergyCount          int                `json:"active_energy_count"`
	ActiveBandwidthCount       int                `json:"active_bandwidth_count"`
	TotalEnergyConsumption     int                `json:"total_energy_consumption"`
	TotalBandwidthConsumption  int                `json:"total_bandwidth_consumption"`
	ActiveEnergyConsumption    int                `json:"active_energy_consumption"`
	ActiveBandwidthConsumption int                `json:"active_bandwidth_consumption"`
	PeriodPricesEnergy         map[string]float64 `json:"period_prices_energy"` // keys are strings "15", "60" etc.
	PeriodPricesBandwidth      map[string]float64 `json:"period_prices_bandwidth"`
	TrxTopUpFee                float64            `json:"trx_top_up_fee"`
	AddressActivationFee       float64            `json:"address_activation_fee"`
	RechargePriceSun           float64            `json:"recharge_price_sun"`
	DailyExpensesAvg           float64            `json:"daily_expenses_avg"`
}

// GetConsumersSummary retrieves summary.
func (c *Client) GetConsumersSummary(ctx context.Context) (*APIResponse[*ConsumersSummary], error) {
	var resp APIResponse[*ConsumersSummary]
	err := c.sendRequest(ctx, "GET", "/api/consumers/summary", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// MassPaymentPeriodParams
type MassPaymentPeriodParams struct {
	ConsumerIDs   []string // or int
	PaymentPeriod int
	AutoRenewal   bool
}

// MassPaymentPeriod updates payment period for multiple consumers.
func (c *Client) MassPaymentPeriod(ctx context.Context, params MassPaymentPeriodParams) (*APIResponse[struct{}], error) {
	// POST /api/consumers/mass/payment-period
	// Samples show `consumer_ids[]="1"`.
	// We need to support array values in form data.
	// postMultipart simple helper doesn't support arrays cleanly, we might need a custom builder here or update helper.
	// Let's do it manually here for correctness.

	path := "/api/consumers/mass/payment-period"
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	// We can use url.Values? No, sample uses form-data.
	// If we use multipart:
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, id := range params.ConsumerIDs {
		writer.WriteField("consumer_ids[]", id)
	}
	writer.WriteField("payment_period", strconv.Itoa(params.PaymentPeriod))
	renewal := "0"
	if params.AutoRenewal {
		renewal = "1"
	}
	writer.WriteField("auto_renewal", renewal)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for k, v := range c.headers {
		req.Header[k] = v
	}

	var resp APIResponse[struct{}]
	_, err = c.Do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// MassTrxParams
type MassTrxParams struct {
	Consumers []string // Ids
	Amount    float64
}

// MassTrx sends TRX to consumers.
func (c *Client) MassTrx(ctx context.Context, params MassTrxParams) (*APIResponse[struct{}], error) {
	// POST /api/consumers/mass/trx

	// Same issue with array consumers[]

	path := "/api/consumers/mass/trx"
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("amount", fmt.Sprintf("%f", params.Amount))
	for _, id := range params.Consumers {
		writer.WriteField("consumers[]", id)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for k, v := range c.headers {
		req.Header[k] = v
	}

	var resp APIResponse[struct{}]
	_, err = c.Do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
