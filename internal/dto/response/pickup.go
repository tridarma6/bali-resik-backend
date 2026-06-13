package response

import (
	"time"

	"github.com/google/uuid"
)

type PickupResponse struct {
	ID            uuid.UUID  `json:"id"`
	WasteType     string     `json:"waste_type"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	Address       string     `json:"address"`
	Status        string     `json:"status"`
	ScheduledDate *time.Time `json:"scheduled_date,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	User          *UserRef   `json:"user,omitempty"`
	Collector     *UserRef   `json:"collector,omitempty"`
}

type UserRef struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone,omitempty"`
}
