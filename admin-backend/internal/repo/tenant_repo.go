package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

type tenantRepo struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *tenantRepo {
	return &tenantRepo{db: db}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

func (r *tenantRepo) FindAll(ctx context.Context) ([]model.Tenant, error) {
	var tenants []model.Tenant
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&tenants).Error; err != nil {
		return nil, fmt.Errorf("failed to find tenants: %w", err)
	}
	return tenants, nil
}

func (r *tenantRepo) FindByID(ctx context.Context, id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find tenant: %w", err)
	}
	return &tenant, nil
}

func (r *tenantRepo) Update(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Save(tenant).Error; err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}
