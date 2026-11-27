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

// RegisterPeminjamanRoutes dikelola oleh Anggota 3.
func RegisterPeminjamanRoutes(api fiber.Router, db *gorm.DB, cfg config.EnvConfig) {
	peminjamanRepo := repository.NewPeminjamanRepository(db)
	barangRepo := repository.NewBarangRepository(db)
	svc := service.NewPeminjamanService(peminjamanRepo, barangRepo)
	h := handler.NewPeminjamanHandler(svc)

	peminjaman := api.Group("/peminjaman", middleware.WajibLogin(cfg.JWTSecret))
	peminjaman.Get("/", h.List)
	peminjaman.Get("/:id", h.Detail)
	peminjaman.Post("/", h.Create)
	peminjaman.Put("/:id", h.Update)
	peminjaman.Delete("/:id", middleware.HanyaPeran("admin"), h.Delete)
}
