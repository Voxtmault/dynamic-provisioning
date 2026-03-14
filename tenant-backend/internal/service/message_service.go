package service

import (
	"context"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
)

type messageService struct {
	repo iface.MessageRepository
}

func NewMessageService(repo iface.MessageRepository) *messageService {
	return &messageService{repo: repo}
}

func (s *messageService) PostMessage(ctx context.Context, req model.CreateMessageRequest) (*model.Message, error) {
	msg := &model.Message{
		HandlerName: req.HandlerName,
		Content:     req.Content,
	}

	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *messageService) GetMessages(ctx context.Context, page, limit int) ([]model.Message, int64, error) {
	return s.repo.FindAll(ctx, page, limit)
}
