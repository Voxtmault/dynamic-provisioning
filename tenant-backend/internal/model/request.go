package model

type CreateMessageRequest struct {
	HandlerName string `json:"handler_name" validate:"required,max=255"`
	Content     string `json:"content" validate:"required,max=1024"`
}
