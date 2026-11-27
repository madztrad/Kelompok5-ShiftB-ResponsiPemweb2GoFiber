package utils

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Sukses bool   `json:"sukses"`
	Pesan  string `json:"pesan,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func Success(c fiber.Ctx, status int, pesan string, data any) error {
	return c.Status(status).JSON(APIResponse{
		Sukses: true,
		Pesan:  pesan,
		Data:   data,
	})
}

func Error(c fiber.Ctx, status int, pesan string) error {
	return c.Status(status).JSON(APIResponse{
		Sukses: false,
		Pesan:  pesan,
	})
}

func ValidationError(c fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(APIResponse{
		Sukses: false,
		Pesan:  "Validasi input gagal",
		Data:   errors,
	})
}
