package trenergy

import (
	"context"
	"fmt"
)

// Partner represents a partner in the structure.
type Partner struct {
	Line  int           `json:"line"`
	Users []PartnerUser `json:"users"`
}

// PartnerUser represents a user in the partner structure.
type PartnerUser struct {
	ID                            int     `json:"id"`
	Name                          string  `json:"name"`
	Photo                         string  `json:"photo"`
	LeaderLevel                   int     `json:"leader_level"`
	LevelName                     string  `json:"level_name"`
	LeaderID                      *int    `json:"leader_id"`
	Stake                         float64 `json:"stake"`
	ActiveStakersCount            int     `json:"active_stakers_count"`
	TotalStakesInStructure        float64 `json:"total_stakes_in_structure"`
	TotalActiveStakersInStructure int     `json:"total_active_stakers_in_structure"`
	TotalPartnersInStructure      int     `json:"total_partners_in_structure"`
}

// GetPartners retrieves the partner structure.
func (c *Client) GetPartners(ctx context.Context, leaderID int) (*APIResponse[[]Partner], error) {
	path := "/api/structure/partners"
	if leaderID != 0 {
		path = fmt.Sprintf("%s?leader=%d", path, leaderID)
	}

	var resp APIResponse[[]Partner]
	err := c.sendRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
