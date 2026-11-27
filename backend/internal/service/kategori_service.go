package service

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
	"siminventaris-backend/internal/repository"
)

type KategoriService struct {
	repo *repository.KategoriRepository
}

func NewKategoriService(repo *repository.KategoriRepository) *KategoriService {
	return &KategoriService{repo: repo}
}

type KategoriInput struct {
	Nama      string `json:"nama" validate:"required,min=2,max=100"`
	Deskripsi string `json:"deskripsi" validate:"max=500"`
}

func (s *KategoriService) List() ([]model.Kategori, error) {
	return s.repo.FindAll()
}

func (s *KategoriService) GetByID(id uint) (*model.Kategori, error) {
	kategori, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusNotFound, Message: "Kategori tidak ditemukan"}
		}
		return nil, err
	}
	return kategori, nil
}

func (s *KategoriService) Create(input KategoriInput) (*model.Kategori, error) {
	nama := strings.TrimSpace(input.Nama)
	if _, err := s.repo.FindByNama(nama); err == nil {
		return nil, &AppError{Code: http.StatusBadRequest, Message: "Nama kategori sudah digunakan"}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	kategori := &model.Kategori{
		Nama:      nama,
		Deskripsi: strings.TrimSpace(input.Deskripsi),
	}
	if err := s.repo.Create(kategori); err != nil {
		return nil, err
	}
	return kategori, nil
}

func (s *KategoriService) Update(id uint, input KategoriInput) (*model.Kategori, error) {
	kategori, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	nama := strings.TrimSpace(input.Nama)
	if existing, err := s.repo.FindByNama(nama); err == nil && existing.ID != id {
		return nil, &AppError{Code: http.StatusBadRequest, Message: "Nama kategori sudah digunakan"}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	kategori.Nama = nama
	kategori.Deskripsi = strings.TrimSpace(input.Deskripsi)
	if err := s.repo.Update(kategori); err != nil {
		return nil, err
	}
	return kategori, nil
}

func (s *KategoriService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}

	count, err := s.repo.CountBarangByKategori(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return &AppError{Code: http.StatusBadRequest, Message: "Kategori tidak dapat dihapus karena masih digunakan oleh barang"}
	}

	return s.repo.Delete(id)
}
