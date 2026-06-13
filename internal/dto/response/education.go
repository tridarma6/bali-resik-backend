package response

import (
	"time"

	"github.com/google/uuid"
)

type ContentResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      *UserRef  `json:"author,omitempty"`
}
