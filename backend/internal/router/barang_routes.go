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

// RegisterBarangRoutes dikelola oleh Anggota 2.
func RegisterBarangRoutes(api fiber.Router, db *gorm.DB, cfg config.EnvConfig) {
	barangRepo := repository.NewBarangRepository(db)
	kategoriRepo := repository.NewKategoriRepository(db)
	svc := service.NewBarangService(barangRepo, kategoriRepo)
	h := handler.NewBarangHandler(svc)

	barang := api.Group("/barang", middleware.WajibLogin(cfg.JWTSecret))
	barang.Get("/", h.List)
	barang.Get("/:id", h.Detail)
	barang.Post("/", middleware.HanyaPeran("admin"), h.Create)
	barang.Put("/:id", middleware.HanyaPeran("admin"), h.Update)
	barang.Delete("/:id", middleware.HanyaPeran("admin"), h.Delete)
}
