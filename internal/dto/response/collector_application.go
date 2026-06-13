package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/domain/models"
)

type CollectorApplicationResponse struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	UserID     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
	AdminNotes string    `json:"admin_notes,omitempty"`
	ReviewedBy *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	User       *UserBriefResponse `json:"user,omitempty"`
}

type UserBriefResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone,omitempty"`
}

func ToCollectorApplicationResponse(app *models.CollectorApplication) *CollectorApplicationResponse {
	resp := &CollectorApplicationResponse{
		ID:         app.ID,
		TenantID:   app.TenantID,
		UserID:     app.UserID,
		Status:     app.Status,
		AdminNotes: app.AdminNotes,
		ReviewedBy: app.ReviewedBy,
		ReviewedAt: app.ReviewedAt,
		CreatedAt:  app.CreatedAt,
	}
	if app.User != nil {
		resp.User = &UserBriefResponse{
			ID:       app.User.ID,
			Email:    app.User.Email,
			FullName: app.User.FullName,
			Phone:    app.User.Phone,
		}
	}
	return resp
}
