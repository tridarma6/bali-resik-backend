package response

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	Type        string     `json:"type"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
