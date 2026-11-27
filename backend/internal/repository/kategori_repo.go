package repository

import (
	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
)

type KategoriRepository struct {
	db *gorm.DB
}

func NewKategoriRepository(db *gorm.DB) *KategoriRepository {
	return &KategoriRepository{db: db}
}

func (r *KategoriRepository) FindAll() ([]model.Kategori, error) {
	var list []model.Kategori
	if err := r.db.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *KategoriRepository) FindByID(id uint) (*model.Kategori, error) {
	var kategori model.Kategori
	if err := r.db.First(&kategori, id).Error; err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *KategoriRepository) FindByNama(nama string) (*model.Kategori, error) {
	var kategori model.Kategori
	if err := r.db.Where("LOWER(nama) = LOWER(?)", nama).First(&kategori).Error; err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *KategoriRepository) CountBarangByKategori(kategoriID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Barang{}).Where("kategori_id = ?", kategoriID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *KategoriRepository) Create(kategori *model.Kategori) error {
	return r.db.Create(kategori).Error
}

func (r *KategoriRepository) Update(kategori *model.Kategori) error {
	return r.db.Save(kategori).Error
}

func (r *KategoriRepository) Delete(id uint) error {
	return r.db.Delete(&model.Kategori{}, id).Error
}
