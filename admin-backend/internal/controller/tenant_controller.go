package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

type TenantController struct {
	service iface.TenantService
}

func NewTenantController(service iface.TenantService) *TenantController {
	return &TenantController{service: service}
}

func (tc *TenantController) RegisterTenant(c echo.Context) error {
	var req model.RegisterTenantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	resp, err := tc.service.RegisterTenant(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to register tenant: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.APIResponse{
		Status:  http.StatusCreated,
		Message: "tenant registered successfully",
		Data:    resp,
	})
}

func (tc *TenantController) ListTenants(c echo.Context) error {
	tenants, err := tc.service.ListTenants(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to list tenants",
		})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Status:  http.StatusOK,
		Message: "tenants retrieved successfully",
		Data:    tenants,
	})
}

func (tc *TenantController) GetTenantProfile(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid tenant id",
		})
	}

	profile, err := tc.service.GetTenantProfile(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, model.APIResponse{
			Status:  http.StatusNotFound,
			Message: "tenant not found",
		})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Status:  http.StatusOK,
		Message: "profile retrieved successfully",
		Data:    profile,
	})
}
