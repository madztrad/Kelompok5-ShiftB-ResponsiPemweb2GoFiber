package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"siminventaris-backend/internal/utils"
)

func WajibLogin(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.Error(c, fiber.StatusUnauthorized, "Token autentikasi tidak ditemukan")
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Error(c, fiber.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Error(c, fiber.StatusUnauthorized, "Klaim token tidak valid")
		}

		uidFloat, ok := claims["user_id"].(float64)
		role, roleOk := claims["role"].(string)
		if !ok || !roleOk {
			return utils.Error(c, fiber.StatusUnauthorized, "Data pengguna di dalam token tidak valid")
		}

		c.Locals("user_id", uint(uidFloat))
		c.Locals("role", role)
		return c.Next()
	}
}
