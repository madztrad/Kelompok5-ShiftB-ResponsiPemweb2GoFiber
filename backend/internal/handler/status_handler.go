package handler

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"siminventaris-backend/internal/utils"
)

type StatusHandler struct {
	db *gorm.DB
}

func NewStatusHandler(db *gorm.DB) *StatusHandler {
	return &StatusHandler{db: db}
}

func (h *StatusHandler) CekStatus(c fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Koneksi database gagal diakses")
	}
	if err := sqlDB.Ping(); err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Database tidak merespons")
	}

	return utils.Success(c, fiber.StatusOK, "API dan database berjalan normal", fiber.Map{
		"status": "ok",
	})
}
