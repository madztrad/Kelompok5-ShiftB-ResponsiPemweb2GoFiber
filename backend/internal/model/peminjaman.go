package model

import "time"

type Peminjaman struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BarangID       uint      `gorm:"not null;index" json:"barang_id"`
	Barang         Barang    `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
	Jumlah         int       `gorm:"not null;default:1" json:"jumlah"`
	TanggalPinjam  string    `gorm:"type:varchar(20);not null" json:"tanggal_pinjam"`
	TanggalKembali string    `gorm:"type:varchar(20)" json:"tanggal_kembali"`
	Status         string    `gorm:"type:varchar(30);not null;default:'dipinjam'" json:"status"`
	Catatan        string    `gorm:"type:text" json:"catatan"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Peminjaman) TableName() string {
	return "peminjaman"
}
