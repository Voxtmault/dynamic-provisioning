package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

type AuthController struct {
	service iface.AuthService
}

func NewAuthController(service iface.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ac *AuthController) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	token, err := ac.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "invalid credentials",
		})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Status:  http.StatusOK,
		Message: "login successful",
		Data:    map[string]string{"token": token},
	})
}
