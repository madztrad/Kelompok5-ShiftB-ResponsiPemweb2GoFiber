package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"siminventaris-backend/internal/config"
	"siminventaris-backend/internal/router"
	"siminventaris-backend/internal/utils"
)

func main() {
	cfg := config.LoadEnv()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	app := fiber.New(fiber.Config{
		// Global Error Handler: cegah detail error internal bocor ke klien
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "Terjadi kesalahan internal server"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				msg = fiberErr.Message
			}
			return utils.Error(c, code, msg)
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	router.SetupRoutes(app, db, cfg)

	log.Printf("Server berjalan di port :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
