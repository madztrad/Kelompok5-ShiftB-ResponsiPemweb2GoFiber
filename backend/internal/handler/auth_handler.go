package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/service"
	"siminventaris-backend/internal/utils"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var input service.RegisterInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}

	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, err := h.svc.Register(input)
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.Success(c, fiber.StatusCreated, "Registrasi berhasil", user)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input service.LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}

	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	token, user, err := h.svc.Login(input)
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "Login berhasil", fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Profile(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	user, err := h.svc.GetProfile(userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "Profil berhasil diambil", user)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return utils.Success(c, fiber.StatusOK, "Logout berhasil, silakan hapus token di sisi klien", nil)
}

func handleServiceError(c fiber.Ctx, err error) error {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		return utils.Error(c, appErr.Code, appErr.Message)
	}
	return utils.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan internal server")
}
