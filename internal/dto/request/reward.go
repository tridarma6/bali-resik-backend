package request

type CreateRewardRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"omitempty,max=500"`
	PointsCost  int    `json:"points_cost" validate:"required,min=1"`
	Stock       *int   `json:"stock,omitempty" validate:"omitempty,min=0"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,max=500"`
}

type UpdateRewardRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=150"`
	Description string `json:"description" validate:"omitempty,max=500"`
	PointsCost  int    `json:"points_cost" validate:"omitempty,min=1"`
	Stock       *int   `json:"stock,omitempty" validate:"omitempty,min=0"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,max=500"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type RedeemRewardRequest struct {
	RewardID string `json:"reward_id" validate:"required,uuid"`
}
