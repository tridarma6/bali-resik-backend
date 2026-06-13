package request

type CreateContentRequest struct {
	Title       string `json:"title" validate:"required,min=5,max=200"`
	Content     string `json:"content" validate:"required,min=20"`
	ContentType string `json:"content_type" validate:"omitempty,oneof=article video infographic"`
	Category    string `json:"category" validate:"required,min=3,max=50"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,max=500"`
}

type UpdateContentRequest struct {
	Title       string `json:"title" validate:"omitempty,min=5,max=200"`
	Content     string `json:"content" validate:"omitempty,min=20"`
	ContentType string `json:"content_type" validate:"omitempty,oneof=article video infographic"`
	Category    string `json:"category" validate:"omitempty,min=3,max=50"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,max=500"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

type ListContentQuery struct {
	Category      string `query:"category"`
	PublishedOnly bool   `query:"published"`
}
