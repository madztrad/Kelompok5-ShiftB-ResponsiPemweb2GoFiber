package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/service"
	"siminventaris-backend/internal/utils"
)

type KategoriHandler struct {
	svc *service.KategoriService
}

func NewKategoriHandler(svc *service.KategoriService) *KategoriHandler {
	return &KategoriHandler{svc: svc}
}

func (h *KategoriHandler) List(c fiber.Ctx) error {
	list, err := h.svc.List()
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Daftar kategori berhasil diambil", list)
}

func (h *KategoriHandler) Detail(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	kategori, err := h.svc.GetByID(id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Detail kategori berhasil diambil", kategori)
}

func (h *KategoriHandler) Create(c fiber.Ctx) error {
	var input service.KategoriInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	kategori, err := h.svc.Create(input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusCreated, "Kategori berhasil ditambahkan", kategori)
}

func (h *KategoriHandler) Update(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	var input service.KategoriInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	kategori, err := h.svc.Update(id, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Kategori berhasil diperbarui", kategori)
}

func (h *KategoriHandler) Delete(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	if err := h.svc.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Kategori berhasil dihapus", nil)
}

func parseUintParam(c fiber.Ctx, name string) (uint, bool) {
	val, err := strconv.ParseUint(c.Params(name), 10, 64)
	if err != nil || val == 0 {
		return 0, false
	}
	return uint(val), true
}
