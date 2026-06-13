package response

import (
	"time"

	"github.com/google/uuid"
)

type RewardResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PointsCost  int       `json:"points_cost"`
	Stock       *int      `json:"stock,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type RewardTransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	Points      int       `json:"points"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserPointsResponse struct {
	TotalPoints int `json:"total_points"`
}
