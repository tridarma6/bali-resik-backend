package response

import (
	"time"

	"github.com/google/uuid"
)

type WasteReportResponse struct {
	ID          uuid.UUID          `json:"id"`
	Latitude    float64            `json:"latitude"`
	Longitude   float64            `json:"longitude"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Severity    string             `json:"severity"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	User        *UserRef           `json:"user,omitempty"`
	Images      []ReportImageResp  `json:"images,omitempty"`
}

type ReportImageResp struct {
	ID       uuid.UUID `json:"id"`
	ImageURL string    `json:"image_url"`
}
