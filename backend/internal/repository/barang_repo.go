package repository

import (
	"strings"

	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
)

type BarangFilter struct {
	Page       int
	Limit      int
	Search     string
	KategoriID uint
	Kondisi    string
}

type BarangRepository struct {
	db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) *BarangRepository {
	return &BarangRepository{db: db}
}

func (r *BarangRepository) FindAll(filter BarangFilter) ([]model.Barang, int64, error) {
	query := r.db.Model(&model.Barang{})

	if s := strings.TrimSpace(strings.ToLower(filter.Search)); s != "" {
		like := "%" + s + "%"
		query = query.Where("LOWER(nama) LIKE ? OR LOWER(kode_barang) LIKE ?", like, like)
	}
	if filter.KategoriID > 0 {
		query = query.Where("kategori_id = ?", filter.KategoriID)
	}
	if k := strings.TrimSpace(filter.Kondisi); k != "" {
		query = query.Where("kondisi = ?", k)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	var items []model.Barang
	if err := query.Preload("Kategori").Order("id DESC").Limit(filter.Limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *BarangRepository) FindByID(id uint) (*model.Barang, error) {
	var barang model.Barang
	if err := r.db.Preload("Kategori").First(&barang, id).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *BarangRepository) FindByKode(kode string) (*model.Barang, error) {
	var barang model.Barang
	if err := r.db.Where("LOWER(kode_barang) = LOWER(?)", kode).First(&barang).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *BarangRepository) CountPeminjamanByBarang(barangID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Peminjaman{}).Where("barang_id = ?", barangID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *BarangRepository) Create(barang *model.Barang) error {
	if err := r.db.Create(barang).Error; err != nil {
		return err
	}
	return r.db.Preload("Kategori").First(barang, barang.ID).Error
}

func (r *BarangRepository) Update(barang *model.Barang) error {
	if err := r.db.Save(barang).Error; err != nil {
		return err
	}
	return r.db.Preload("Kategori").First(barang, barang.ID).Error
}

func (r *BarangRepository) Delete(id uint) error {
	return r.db.Delete(&model.Barang{}, id).Error
}
