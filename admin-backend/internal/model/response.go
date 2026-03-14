package model

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RegisterTenantResponse struct {
	TenantID       uint   `json:"tenant_id"`
	Subdomain      string `json:"subdomain"`
	PhotoUploadURL string `json:"photo_upload_url"`
}

type TenantListItem struct {
	ID                      uint   `json:"id"`
	Name                    string `json:"name"`
	Subdomain               string `json:"subdomain"`
	Status                  string `json:"status"`
	BackendContainerStatus  string `json:"backend_container_status"`
	FrontendContainerStatus string `json:"frontend_container_status"`
}

type TenantProfileResponse struct {
	TenantID     string   `json:"tenant_id"`
	AppName      string   `json:"app_name"`
	AppPhotoKey  string   `json:"app_photo_key"`
	ColorPalette []string `json:"color_palette"`
}
