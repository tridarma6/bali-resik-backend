package usecase

import (
	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
)

type EducationUseCase interface {
	CreateContent(tenantID, userID uuid.UUID, req *request.CreateContentRequest) (*response.ContentResponse, error)
	GetContent(contentID uuid.UUID) (*response.ContentResponse, error)
	ListContents(tenantID uuid.UUID, query request.ListContentQuery) ([]response.ContentResponse, error)
	UpdateContent(contentID uuid.UUID, req *request.UpdateContentRequest) (*response.ContentResponse, error)
	DeleteContent(contentID uuid.UUID) error
}
