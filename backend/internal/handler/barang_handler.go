package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/repository"
	"siminventaris-backend/internal/service"
	"siminventaris-backend/internal/utils"
)

type BarangHandler struct {
	svc *service.BarangService
}

func NewBarangHandler(svc *service.BarangService) *BarangHandler {
	return &BarangHandler{svc: svc}
}

func (h *BarangHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	katID, _ := strconv.ParseUint(c.Query("kategori_id", "0"), 10, 64)

	result, err := h.svc.List(repository.BarangFilter{
		Page:       page,
		Limit:      limit,
		Search:     c.Query("search"),
		KategoriID: uint(katID),
		Kondisi:    c.Query("kondisi"),
	})
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Daftar barang berhasil diambil", result)
}

func (h *BarangHandler) Detail(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID barang tidak valid")
	}

	barang, err := h.svc.GetByID(id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Detail barang berhasil diambil", barang)
}

func (h *BarangHandler) Create(c fiber.Ctx) error {
	var input service.BarangInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	barang, err := h.svc.Create(input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusCreated, "Barang berhasil ditambahkan", barang)
}

func (h *BarangHandler) Update(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID barang tidak valid")
	}

	var input service.BarangInput
	if err := c.Bind().Body(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Format request body tidak valid")
	}
	if errs := utils.ValidateStruct(input); errs != nil {
		return utils.ValidationError(c, errs)
	}

	barang, err := h.svc.Update(id, input)
	if err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Barang berhasil diperbarui", barang)
}

func (h *BarangHandler) Delete(c fiber.Ctx) error {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return utils.Error(c, fiber.StatusBadRequest, "ID barang tidak valid")
	}

	if err := h.svc.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return utils.Success(c, fiber.StatusOK, "Barang berhasil dihapus", nil)
}
