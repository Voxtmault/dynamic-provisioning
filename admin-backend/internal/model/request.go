package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterTenantRequest struct {
	Name         string   `json:"name" validate:"required,max=255"`
	ColorPalette []string `json:"color_palette" validate:"required,min=1"`
}
