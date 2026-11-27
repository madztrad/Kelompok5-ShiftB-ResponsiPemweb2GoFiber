package router

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"siminventaris-backend/internal/config"
	"siminventaris-backend/internal/handler"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg config.EnvConfig) {
	api := app.Group("/api")

	statusHandler := handler.NewStatusHandler(db)
	api.Get("/status", statusHandler.CekStatus)

	// Daftarkan route modular per fitur (1 file = 1 anggota/modul)
	RegisterAuthRoutes(api, db, cfg)
	RegisterKategoriRoutes(api, db, cfg)
	RegisterBarangRoutes(api, db, cfg)
	RegisterPeminjamanRoutes(api, db, cfg)
}
