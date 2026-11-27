package router

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"siminventaris-backend/internal/config"
	"siminventaris-backend/internal/handler"
	"siminventaris-backend/internal/middleware"
	"siminventaris-backend/internal/repository"
	"siminventaris-backend/internal/service"
)

// RegisterAuthRoutes dikelola oleh Anggota 1.
func RegisterAuthRoutes(api fiber.Router, db *gorm.DB, cfg config.EnvConfig) {
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authSvc)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/profile", middleware.WajibLogin(cfg.JWTSecret), authHandler.Profile)
	auth.Post("/logout", middleware.WajibLogin(cfg.JWTSecret), authHandler.Logout)
}
