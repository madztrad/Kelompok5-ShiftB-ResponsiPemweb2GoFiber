package config

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"siminventaris-backend/internal/model"
)

func ConnectDB(cfg EnvConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Kategori{},
		&model.Barang{},
		&model.Peminjaman{},
	); err != nil {
		return nil, err
	}

	seedDemoUsers(db)
	return db, nil
}

func seedDemoUsers(db *gorm.DB) {
	demoAccounts := []struct {
		Nama     string
		Email    string
		Password string
		Role     string
	}{
		{"Admin Lab", "admin@lab.test", "Admin123!", "admin"},
		{"Mahasiswa Demo", "member@lab.test", "Member123!", "member"},
	}

	for _, acc := range demoAccounts {
		var existing model.User
		err := db.Where("email = ?", acc.Email).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hash, hashErr := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
			if hashErr != nil {
				log.Printf("Gagal hash password seeder %s: %v\n", acc.Email, hashErr)
				continue
			}
			user := model.User{
				Nama:     acc.Nama,
				Email:    acc.Email,
				Password: string(hash),
				Role:     acc.Role,
			}
			if createErr := db.Create(&user).Error; createErr != nil {
				log.Printf("Gagal membuat akun demo %s: %v\n", acc.Email, createErr)
			}
		}
	}
}
