package config

import (
	"fmt"
	"os"
)

type Config struct {
	OpenBaoAddr     string
	OpenBaoRoleID   string
	OpenBaoSecretID string
	TenantID        string
	Env             string
	AdminBackendURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		OpenBaoAddr:     os.Getenv("OPENBAO_ADDR"),
		OpenBaoRoleID:   os.Getenv("OPENBAO_ROLE_ID"),
		OpenBaoSecretID: os.Getenv("OPENBAO_SECRET_ID"),
		TenantID:        os.Getenv("TENANT_ID"),
		Env:             os.Getenv("ENV"),
		AdminBackendURL: os.Getenv("ADMIN_BACKEND_URL"),
	}

	required := map[string]string{
		"OPENBAO_ADDR":      cfg.OpenBaoAddr,
		"OPENBAO_ROLE_ID":   cfg.OpenBaoRoleID,
		"OPENBAO_SECRET_ID": cfg.OpenBaoSecretID,
		"TENANT_ID":         cfg.TenantID,
		"ENV":               cfg.Env,
		"ADMIN_BACKEND_URL": cfg.AdminBackendURL,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return cfg, nil
}
