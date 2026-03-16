package router

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/controller"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/validator"
)

func Setup(
	e *echo.Echo,
	authCtrl *controller.AuthController,
	tenantCtrl *controller.TenantController,
	jwtSecret string,
) {
	// Validator
	e.Validator = validator.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api := e.Group("/api")

	// Public routes
	api.POST("/auth/login", authCtrl.Login)
	api.GET("/tenant/:id/profile", tenantCtrl.GetTenantProfile)

	// Protected routes
	protected := api.Group("")
	protected.Use(jwtMiddleware(jwtSecret))
	protected.POST("/tenants", tenantCtrl.RegisterTenant)
	protected.GET("/tenants", tenantCtrl.ListTenants)
	protected.PUT("/tenants/:id/restart", tenantCtrl.RestartTenant)
}

func jwtMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, model.APIResponse{
					Status:  http.StatusUnauthorized,
					Message: "missing authorization header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, model.APIResponse{
					Status:  http.StatusUnauthorized,
					Message: "invalid authorization format",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, model.APIResponse{
					Status:  http.StatusUnauthorized,
					Message: "invalid or expired token",
				})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("user_id", claims["user_id"])
				c.Set("email", claims["email"])
			}

			return next(c)
		}
	}
}
