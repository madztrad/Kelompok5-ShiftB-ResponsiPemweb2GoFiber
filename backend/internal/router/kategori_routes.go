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

// RegisterKategoriRoutes dikelola oleh Anggota 1.
func RegisterKategoriRoutes(api fiber.Router, db *gorm.DB, cfg config.EnvConfig) {
	repo := repository.NewKategoriRepository(db)
	svc := service.NewKategoriService(repo)
	h := handler.NewKategoriHandler(svc)

	kategori := api.Group("/kategori", middleware.WajibLogin(cfg.JWTSecret))
	kategori.Get("/", h.List)
	kategori.Get("/:id", h.Detail)
	kategori.Post("/", middleware.HanyaPeran("admin"), h.Create)
	kategori.Put("/:id", middleware.HanyaPeran("admin"), h.Update)
	kategori.Delete("/:id", middleware.HanyaPeran("admin"), h.Delete)
}
