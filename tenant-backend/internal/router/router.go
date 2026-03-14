package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/controller"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/validator"
)

func Setup(
	e *echo.Echo,
	msgCtrl *controller.MessageController,
	profileCtrl *controller.ProfileController,
) {
	// Validator
	e.Validator = validator.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	api := e.Group("/api")
	api.POST("/messages", msgCtrl.PostMessage)
	api.GET("/messages", msgCtrl.GetMessages)
	api.GET("/profile", profileCtrl.GetProfile)
}
