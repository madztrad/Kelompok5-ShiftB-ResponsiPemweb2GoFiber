package model

import "time"

type Barang struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	KodeBarang string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"kode_barang"`
	Nama       string    `gorm:"type:varchar(150);not null" json:"nama"`
	Stok       int       `gorm:"not null;default:0" json:"stok"`
	Kondisi    string    `gorm:"type:varchar(30);not null;default:'baik'" json:"kondisi"`
	KategoriID uint      `gorm:"not null;index" json:"kategori_id"`
	Kategori   Kategori  `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Barang) TableName() string {
	return "barang"
}
