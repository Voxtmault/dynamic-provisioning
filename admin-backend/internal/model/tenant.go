package model

import (
	"time"

	"github.com/lib/pq"
)

type Tenant struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Subdomain    string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"subdomain"`
	AppPhotoKey  string         `gorm:"type:varchar(512)" json:"app_photo_key,omitempty"`
	ColorPalette pq.StringArray `gorm:"type:text[]" json:"color_palette"`
	Status       string         `gorm:"type:varchar(50);not null;default:'provisioning'" json:"status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

const (
	TenantStatusProvisioning = "provisioning"
	TenantStatusActive       = "active"
	TenantStatusStopped      = "stopped"
	TenantStatusError        = "error"
)
