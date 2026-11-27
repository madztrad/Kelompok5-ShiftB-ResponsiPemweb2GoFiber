package model

import "time"

type Kategori struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"nama"`
	Deskripsi string    `gorm:"type:text" json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Kategori) TableName() string {
	return "kategori"
}
