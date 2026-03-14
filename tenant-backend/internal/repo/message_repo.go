package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
)

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepo {
	return &messageRepo{db: db}
}

func (r *messageRepo) Create(ctx context.Context, msg *model.Message) error {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *messageRepo) FindAll(ctx context.Context, page, limit int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Message{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find messages: %w", err)
	}

	return messages, total, nil
}
