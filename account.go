package trenergy

import (
	"context"
)

// AccountInfo represents user account information.
type AccountInfo struct {
	Name                  string                `json:"name"`
	Email                 string                `json:"email"`
	HasPassword           bool                  `json:"has_password"`
	HasTelegram           bool                  `json:"has_telegram"`
	Lang                  string                `json:"lang"`
	TheCode               string                `json:"the_code"`
	InvitationCode        *string               `json:"invitation_code"`
	CreditLimit           float64               `json:"credit_limit"`
	LeaderName            *string               `json:"leader_name"`
	LeaderLevel           int                   `json:"leader_level"`
	RefEnabled            bool                  `json:"ref_enabled"`
	IsBanned              bool                  `json:"is_banned"`
	BalanceRestricted     bool                  `json:"balance_restricted"`
	Balance               float64               `json:"balance"`
	EnergyBalance         float64               `json:"energy_balance"`
	Photo                 string                `json:"photo"`
	StakesSum             float64               `json:"stakes_sum"`
	StakesProfit          float64               `json:"stakes_profit"`
	AvailableToUnstakeSum float64               `json:"available_to_unstake_sum"`
	ActiveStakersCount    int                   `json:"active_stakers_count"`
	Subscription          *Subscription         `json:"subscription"`
	TwoFA                 bool                  `json:"2fa"`
	Onboarding            int                   `json:"onboarding"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
	DeletionAt            *string               `json:"deletion_at"`
	Reinvestment          *Reinvestment         `json:"reinvestment"`
	NotificationSettings  []NotificationSetting `json:"notification_settings"`
}

type Subscription struct {
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

type Reinvestment struct {
	Wallet    bool   `json:"wallet"`
	Balance   bool   `json:"balance"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NotificationSetting struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

// GetAccountInfo retrieves the current user's account information.
func (c *Client) GetAccountInfo(ctx context.Context) (*APIResponse[*AccountInfo], error) {
	var resp APIResponse[*AccountInfo]
	err := c.sendRequest(ctx, "GET", "/api/account", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// TopUpInfo represents info for topping up account.
type TopUpInfo struct {
	Address  string `json:"address"`
	QRCode   string `json:"qr_code"`
	TimeLeft int    `json:"time_left"`
}

// GetTopUpInfo retrieves information for account top-up.
func (c *Client) GetTopUpInfo(ctx context.Context) (*APIResponse[*TopUpInfo], error) {
	var resp APIResponse[*TopUpInfo]
	err := c.sendRequest(ctx, "GET", "/api/account/top-up", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
