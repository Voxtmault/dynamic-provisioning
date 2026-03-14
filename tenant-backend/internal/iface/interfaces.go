package iface

import (
	"context"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
	FindAll(ctx context.Context, page, limit int) ([]model.Message, int64, error)
}

type MessageService interface {
	PostMessage(ctx context.Context, req model.CreateMessageRequest) (*model.Message, error)
	GetMessages(ctx context.Context, page, limit int) ([]model.Message, int64, error)
}

type ProfileService interface {
	GetProfile(ctx context.Context) (*model.AppProfile, error)
}
