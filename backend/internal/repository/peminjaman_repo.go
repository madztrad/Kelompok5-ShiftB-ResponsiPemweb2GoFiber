package repository

import (
	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
)

type PeminjamanRepository struct {
	db *gorm.DB
}

func NewPeminjamanRepository(db *gorm.DB) *PeminjamanRepository {
	return &PeminjamanRepository{db: db}
}

func (r *PeminjamanRepository) DB() *gorm.DB {
	return r.db
}

func (r *PeminjamanRepository) FindAll(userID uint, isAdmin bool) ([]model.Peminjaman, error) {
	query := r.db.Preload("User").Preload("Barang").Preload("Barang.Kategori").Order("id DESC")
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}
	list := make([]model.Peminjaman, 0)
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PeminjamanRepository) FindByID(id uint) (*model.Peminjaman, error) {
	var p model.Peminjaman
	if err := r.db.Preload("User").Preload("Barang").Preload("Barang.Kategori").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
