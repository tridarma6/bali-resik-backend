package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/domain/models"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type EducationUseCaseImpl struct {
	eduRepo repo.EducationRepository
	log     *logrus.Logger
}

func NewEducationUseCase(
	eduRepo repo.EducationRepository,
	log *logrus.Logger,
) *EducationUseCaseImpl {
	return &EducationUseCaseImpl{
		eduRepo: eduRepo,
		log:     log,
	}
}

func (uc *EducationUseCaseImpl) CreateContent(tenantID, userID uuid.UUID, req *request.CreateContentRequest) (*response.ContentResponse, error) {
	contentType := req.ContentType
	if contentType == "" {
		contentType = "article"
	}

	content := &models.EducationalContent{
		TenantID:    tenantID,
		Title:       req.Title,
		Content:     req.Content,
		ContentType: contentType,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		CreatedBy:   userID,
	}

	if err := uc.eduRepo.Create(content); err != nil {
		uc.log.WithError(err).Error("failed to create educational content")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to create content", err)
	}

	return uc.toContentResponse(content), nil
}

func (uc *EducationUseCaseImpl) GetContent(contentID uuid.UUID) (*response.ContentResponse, error) {
	content, err := uc.eduRepo.FindByID(contentID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Content not found", errors.ErrNotFound)
	}

	return uc.toContentResponse(content), nil
}

func (uc *EducationUseCaseImpl) ListContents(tenantID uuid.UUID, query request.ListContentQuery) ([]response.ContentResponse, error) {
	contents, err := uc.eduRepo.FindByTenant(tenantID, query.Category, query.PublishedOnly)
	if err != nil {
		uc.log.WithError(err).Error("failed to list educational content")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to list content", err)
	}

	result := make([]response.ContentResponse, 0, len(contents))
	for _, c := range contents {
		result = append(result, *uc.toContentResponse(&c))
	}

	return result, nil
}

func (uc *EducationUseCaseImpl) UpdateContent(contentID uuid.UUID, req *request.UpdateContentRequest) (*response.ContentResponse, error) {
	content, err := uc.eduRepo.FindByID(contentID)
	if err != nil {
		return nil, errors.NewDomainError("NOT_FOUND", "Content not found", errors.ErrNotFound)
	}

	if req.Title != "" {
		content.Title = req.Title
	}
	if req.Content != "" {
		content.Content = req.Content
	}
	if req.ContentType != "" {
		content.ContentType = req.ContentType
	}
	if req.Category != "" {
		content.Category = req.Category
	}
	if req.ImageURL != "" {
		content.ImageURL = req.ImageURL
	}
	if req.IsPublished != nil {
		content.IsPublished = *req.IsPublished
	}

	if err := uc.eduRepo.Update(content); err != nil {
		uc.log.WithError(err).Error("failed to update educational content")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to update content", err)
	}

	return uc.toContentResponse(content), nil
}

func (uc *EducationUseCaseImpl) DeleteContent(contentID uuid.UUID) error {
	if err := uc.eduRepo.Delete(contentID); err != nil {
		return errors.NewDomainError("NOT_FOUND", "Content not found", errors.ErrNotFound)
	}
	return nil
}

func (uc *EducationUseCaseImpl) toContentResponse(c *models.EducationalContent) *response.ContentResponse {
	resp := &response.ContentResponse{
		ID:          c.ID,
		Title:       c.Title,
		Content:     c.Content,
		ContentType: c.ContentType,
		Category:    c.Category,
		ImageURL:    c.ImageURL,
		IsPublished: c.IsPublished,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	if c.Author != nil {
		resp.Author = &response.UserRef{
			ID:       c.Author.ID,
			FullName: c.Author.FullName,
			Email:    c.Author.Email,
		}
	}

	return resp
}
