package main

import (
	"github.com/iv-tunate/bizpapertrail/handlers"
	middlewares "github.com/iv-tunate/bizpapertrail/internal/middleware"
	"github.com/labstack/echo"
	 middleware "github.com/labstack/echo/middleware"
)

func registerRoutes(e *echo.Echo, h *handlers.Handler) {
	api := e.Group("/api")

	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge: 300,
		ExposeHeaders:    []string{"Link"},
	}))
	
	protected := api.Group("")
	protected.Use(middlewares.JWTMiddleware)

	user := protected.Group("/user")
	user.POST("/user", h.RegisterUser)
	user.PUT("/verify_email", h.VerifyUser)

	auth := protected.Group("/auth")
	auth.POST("/refresh", h.RefreshJwtToken)
}