package request

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
	FullName    string `json:"full_name" validate:"required,min=2,max=150"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=20"`
	TenantSlug  string `json:"tenant_slug" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CreateTenantRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Slug       string `json:"slug" validate:"required,min=2,max=100"`
	RegionType string `json:"region_type" validate:"required,oneof=kota kabupaten"`
}

type CreateAdminRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FullName  string `json:"full_name" validate:"required,min=2,max=150"`
	TenantID  string `json:"tenant_id" validate:"required,uuid"`
}
