package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/utils"
)

// HanyaPeran WAJIB diletakkan SEBELUM handler pada definisi route Fiber v3.
func HanyaPeran(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role == "" || !slices.Contains(roles, role) {
			return utils.Error(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk aksi ini")
		}
		return c.Next()
	}
}
