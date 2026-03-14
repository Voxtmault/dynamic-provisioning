package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
)

type MessageController struct {
	service iface.MessageService
}

func NewMessageController(service iface.MessageService) *MessageController {
	return &MessageController{service: service}
}

func (mc *MessageController) PostMessage(c echo.Context) error {
	var req model.CreateMessageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	msg, err := mc.service.PostMessage(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to post message",
		})
	}

	return c.JSON(http.StatusCreated, model.APIResponse{
		Status:  http.StatusCreated,
		Message: "message posted successfully",
		Data:    msg,
	})
}

func (mc *MessageController) GetMessages(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	messages, total, err := mc.service.GetMessages(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve messages",
		})
	}

	return c.JSON(http.StatusOK, model.PaginatedResponse{
		Status:  http.StatusOK,
		Message: "messages retrieved successfully",
		Data:    messages,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}
