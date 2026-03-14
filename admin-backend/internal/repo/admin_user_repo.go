package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

type adminUserRepo struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *adminUserRepo {
	return &adminUserRepo{db: db}
}

func (r *adminUserRepo) Create(ctx context.Context, user *model.AdminUser) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	return nil
}

func (r *adminUserRepo) FindByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	var user model.AdminUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find admin user: %w", err)
	}
	return &user, nil
}
