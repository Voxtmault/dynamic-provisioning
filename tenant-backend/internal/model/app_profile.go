package model

type AppProfile struct {
	TenantID     string   `json:"tenant_id"`
	AppName      string   `json:"app_name"`
	AppPhotoURL  string   `json:"app_photo_url"`
	AppPhotoKey  string   `json:"app_photo_key,omitempty"` // S3 object key, not exposed to frontend
	ColorPalette []string `json:"color_palette"`
}
