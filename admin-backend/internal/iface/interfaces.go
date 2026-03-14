package iface

import (
	"context"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

// Repositories

type AdminUserRepository interface {
	Create(ctx context.Context, user *model.AdminUser) error
	FindByEmail(ctx context.Context, email string) (*model.AdminUser, error)
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	FindAll(ctx context.Context) ([]model.Tenant, error)
	FindByID(ctx context.Context, id uint) (*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
}

// Services

type AuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (string, error)
	SeedAdmin(ctx context.Context, email, password string) error
}

type TenantService interface {
	RegisterTenant(ctx context.Context, req model.RegisterTenantRequest) (*model.RegisterTenantResponse, error)
	ListTenants(ctx context.Context) ([]model.TenantListItem, error)
	GetTenantProfile(ctx context.Context, tenantID uint) (*model.TenantProfileResponse, error)
}

type ProvisioningService interface {
	ProvisionTenant(ctx context.Context, tenant *model.Tenant) (roleID, secretID string, err error)
	SpinUpContainers(ctx context.Context, tenant *model.Tenant, roleID, secretID string) error
}
