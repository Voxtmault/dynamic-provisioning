package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
)

type ProfileController struct {
	service iface.ProfileService
}

func NewProfileController(service iface.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

func (pc *ProfileController) GetProfile(c echo.Context) error {
	profile, err := pc.service.GetProfile(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve app profile",
		})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Status:  http.StatusOK,
		Message: "profile retrieved successfully",
		Data:    profile,
	})
}
