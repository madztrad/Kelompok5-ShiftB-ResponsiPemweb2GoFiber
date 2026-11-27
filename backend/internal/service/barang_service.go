package service

import (
	"errors"
	"math"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
	"siminventaris-backend/internal/repository"
)

type BarangService struct {
	repo         *repository.BarangRepository
	kategoriRepo *repository.KategoriRepository
}

func NewBarangService(repo *repository.BarangRepository, kategoriRepo *repository.KategoriRepository) *BarangService {
	return &BarangService{repo: repo, kategoriRepo: kategoriRepo}
}

type BarangInput struct {
	KodeBarang string `json:"kode_barang" validate:"required,min=2,max=50"`
	Nama       string `json:"nama" validate:"required,min=2,max=150"`
	Stok       int    `json:"stok" validate:"gte=0"`
	Kondisi    string `json:"kondisi" validate:"required,oneof=baik rusak_ringan rusak_berat"`
	KategoriID uint   `json:"kategori_id" validate:"required,gt=0"`
}

type PaginatedBarang struct {
	Items []model.Barang `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalData    int64 `json:"total_data"`
	TotalHalaman int   `json:"total_halaman"`
}

func (s *BarangService) List(filter repository.BarangFilter) (*PaginatedBarang, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}

	items, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalHalaman := int(math.Ceil(float64(total) / float64(filter.Limit)))
	if totalHalaman == 0 {
		totalHalaman = 1
	}

	return &PaginatedBarang{
		Items: items,
		Meta: PaginationMeta{
			Page:         filter.Page,
			Limit:        filter.Limit,
			TotalData:    total,
			TotalHalaman: totalHalaman,
		},
	}, nil
}

func (s *BarangService) GetByID(id uint) (*model.Barang, error) {
	barang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusNotFound, Message: "Barang tidak ditemukan"}
		}
		return nil, err
	}
	return barang, nil
}

func (s *BarangService) Create(input BarangInput) (*model.Barang, error) {
	if _, err := s.kategoriRepo.FindByID(input.KategoriID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusBadRequest, Message: "Kategori tidak ditemukan"}
		}
		return nil, err
	}

	kode := strings.ToUpper(strings.TrimSpace(input.KodeBarang))
	if _, err := s.repo.FindByKode(kode); err == nil {
		return nil, &AppError{Code: http.StatusBadRequest, Message: "Kode barang sudah digunakan"}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	barang := &model.Barang{
		KodeBarang: kode,
		Nama:       strings.TrimSpace(input.Nama),
		Stok:       input.Stok,
		Kondisi:    input.Kondisi,
		KategoriID: input.KategoriID,
	}

	if err := s.repo.Create(barang); err != nil {
		return nil, err
	}
	return barang, nil
}

func (s *BarangService) Update(id uint, input BarangInput) (*model.Barang, error) {
	barang, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if _, err := s.kategoriRepo.FindByID(input.KategoriID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusBadRequest, Message: "Kategori tidak ditemukan"}
		}
		return nil, err
	}

	kode := strings.ToUpper(strings.TrimSpace(input.KodeBarang))
	if existing, err := s.repo.FindByKode(kode); err == nil && existing.ID != id {
		return nil, &AppError{Code: http.StatusBadRequest, Message: "Kode barang sudah digunakan"}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	barang.KodeBarang = kode
	barang.Nama = strings.TrimSpace(input.Nama)
	barang.Stok = input.Stok
	barang.Kondisi = input.Kondisi
	barang.KategoriID = input.KategoriID

	if err := s.repo.Update(barang); err != nil {
		return nil, err
	}
	return barang, nil
}

func (s *BarangService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}

	count, err := s.repo.CountPeminjamanByBarang(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return &AppError{Code: http.StatusBadRequest, Message: "Barang tidak dapat dihapus karena memiliki riwayat peminjaman"}
	}

	return s.repo.Delete(id)
}
