package handler

import (
	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/service"
	"siminventaris-backend/internal/utils"
)

type PeminjamanHandler struct {
	svc *service.PeminjamanService
}

func NewPeminjamanHandler(svc *service.PeminjamanService) *PeminjamanHandler {
	return &PeminjamanHandler{svc: svc}
}

func (h *PeminjamanHandler) List(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(uint)
	role, _ := c.Locals("role").(string)

	list, err := h.svc.List(userID, role)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Daftar peminjaman berhasil diambil", list)
}

func (h *PeminjamanHandler) Detail(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID peminjaman tidak valid")
	}

	userID, _ := c.Locals("user_id").(uint)
	role, _ := c.Locals("role").(string)

	p, err := h.svc.GetByID(id, userID, role)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Detail peminjaman berhasil diambil", p)
}

func (h *PeminjamanHandler) Create(c fiber.Ctx) error {
	var input service.CreatePeminjamanInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	userID, _ := c.Locals("user_id").(uint)
	p, err := h.svc.Create(userID, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusCreated, "Peminjaman berhasil diajukan", p)
}

func (h *PeminjamanHandler) Update(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID peminjaman tidak valid")
	}

	var input service.UpdatePeminjamanInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	userID, _ := c.Locals("user_id").(uint)
	role, _ := c.Locals("role").(string)

	p, err := h.svc.Update(id, userID, role, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Peminjaman berhasil diperbarui", p)
}

func (h *PeminjamanHandler) Delete(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID peminjaman tidak valid")
	}

	if err := h.svc.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Peminjaman berhasil dihapus", nil)
}
