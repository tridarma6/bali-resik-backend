package request

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=150"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=20"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}
