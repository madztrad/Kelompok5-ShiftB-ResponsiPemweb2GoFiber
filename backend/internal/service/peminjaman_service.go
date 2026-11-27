package service

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
	"siminventaris-backend/internal/repository"
)

type PeminjamanService struct {
	repo       *repository.PeminjamanRepository
	barangRepo *repository.BarangRepository
}

func NewPeminjamanService(repo *repository.PeminjamanRepository, barangRepo *repository.BarangRepository) *PeminjamanService {
	return &PeminjamanService{repo: repo, barangRepo: barangRepo}
}

type CreatePeminjamanInput struct {
	BarangID       uint   `json:"barang_id" validate:"required,gt=0"`
	Jumlah         int    `json:"jumlah" validate:"required,gte=1"`
	TanggalPinjam  string `json:"tanggal_pinjam" validate:"required,min=8,max=20"`
	TanggalKembali string `json:"tanggal_kembali" validate:"max=20"`
	Catatan        string `json:"catatan" validate:"max=500"`
}

type UpdatePeminjamanInput struct {
	Jumlah         *int   `json:"jumlah" validate:"omitempty,gte=1"`
	TanggalPinjam  string `json:"tanggal_pinjam" validate:"omitempty,min=8,max=20"`
	TanggalKembali string `json:"tanggal_kembali" validate:"max=20"`
	Status         string `json:"status" validate:"omitempty,oneof=dipinjam dikembalikan ditolak"`
	Catatan        string `json:"catatan" validate:"max=500"`
}

func (s *PeminjamanService) List(userID uint, role string) ([]model.Peminjaman, error) {
	return s.repo.FindAll(userID, role == "admin")
}

func (s *PeminjamanService) GetByID(id uint, userID uint, role string) (*model.Peminjaman, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusNotFound, Message: "Data peminjaman tidak ditemukan"}
		}
		return nil, err
	}
	if role != "admin" && p.UserID != userID {
		return nil, &AppError{Code: http.StatusForbidden, Message: "Akses ditolak: Anda hanya dapat melihat peminjaman milik sendiri"}
	}
	return p, nil
}

func (s *PeminjamanService) Create(userID uint, input CreatePeminjamanInput) (*model.Peminjaman, error) {
	var createdID uint

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var barang model.Barang
		if err := tx.First(&barang, input.BarangID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &AppError{Code: http.StatusNotFound, Message: "Barang tidak ditemukan"}
			}
			return err
		}

		if barang.Stok < input.Jumlah {
			return &AppError{Code: http.StatusBadRequest, Message: "Stok barang tidak mencukupi"}
		}

		barang.Stok -= input.Jumlah
		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

		p := model.Peminjaman{
			UserID:         userID,
			BarangID:       input.BarangID,
			Jumlah:         input.Jumlah,
			TanggalPinjam:  strings.TrimSpace(input.TanggalPinjam),
			TanggalKembali: strings.TrimSpace(input.TanggalKembali),
			Status:         "dipinjam",
			Catatan:        strings.TrimSpace(input.Catatan),
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		createdID = p.ID
		return nil
	})

	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(createdID)
}

func (s *PeminjamanService) Update(id uint, userID uint, role string, input UpdatePeminjamanInput) (*model.Peminjaman, error) {
	p, err := s.GetByID(id, userID, role)
	if err != nil {
		return nil, err
	}

	if role != "admin" && input.Status != "" && input.Status != p.Status {
		return nil, &AppError{Code: http.StatusForbidden, Message: "Hanya admin yang dapat mengubah status peminjaman"}
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var barang model.Barang
		if err := tx.First(&barang, p.BarangID).Error; err != nil {
			return err
		}

		oldStatus := p.Status
		oldJumlah := p.Jumlah

		if input.Jumlah != nil && *input.Jumlah != oldJumlah {
			if oldStatus == "dipinjam" {
				diff := *input.Jumlah - oldJumlah
				if barang.Stok < diff {
					return &AppError{Code: http.StatusBadRequest, Message: "Stok barang tidak mencukupi untuk penambahan jumlah"}
				}
				barang.Stok -= diff
			}
			p.Jumlah = *input.Jumlah
		}

		if input.TanggalPinjam != "" {
			p.TanggalPinjam = strings.TrimSpace(input.TanggalPinjam)
		}
		if input.TanggalKembali != "" {
			p.TanggalKembali = strings.TrimSpace(input.TanggalKembali)
		}
		if input.Catatan != "" {
			p.Catatan = strings.TrimSpace(input.Catatan)
		}

		if role == "admin" && input.Status != "" && input.Status != oldStatus {
			if oldStatus == "dipinjam" && (input.Status == "dikembalikan" || input.Status == "ditolak") {
				barang.Stok += p.Jumlah
			} else if (oldStatus == "dikembalikan" || oldStatus == "ditolak") && input.Status == "dipinjam" {
				if barang.Stok < p.Jumlah {
					return &AppError{Code: http.StatusBadRequest, Message: "Stok barang tidak mencukupi untuk mengaktifkan kembali peminjaman"}
				}
				barang.Stok -= p.Jumlah
			}
			p.Status = input.Status
		}

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}
		return tx.Save(p).Error
	})

	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *PeminjamanService) Delete(id uint) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var p model.Peminjaman
		if err := tx.First(&p, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &AppError{Code: http.StatusNotFound, Message: "Data peminjaman tidak ditemukan"}
			}
			return err
		}

		if p.Status == "dipinjam" {
			var barang model.Barang
			if err := tx.First(&barang, p.BarangID).Error; err == nil {
				barang.Stok += p.Jumlah
				if err := tx.Save(&barang).Error; err != nil {
					return err
				}
			}
		}

		return tx.Delete(&model.Peminjaman{}, id).Error
	})
}
