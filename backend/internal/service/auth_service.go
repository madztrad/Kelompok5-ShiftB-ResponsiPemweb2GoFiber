package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siminventaris-backend/internal/config"
	"siminventaris-backend/internal/model"
	"siminventaris-backend/internal/repository"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

type AuthService struct {
	repo *repository.UserRepository
	cfg  config.EnvConfig
}

func NewAuthService(repo *repository.UserRepository, cfg config.EnvConfig) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

type RegisterInput struct {
	Nama     string `json:"nama" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *AuthService) Register(input RegisterInput) (*model.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, &AppError{Code: http.StatusBadRequest, Message: "Email sudah terdaftar"}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Registrasi publik wajib dipaksa menjadi role 'member'
	user := &model.User{
		Nama:     strings.TrimSpace(input.Nama),
		Email:    email,
		Password: string(hash),
		Role:     "member",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(input LoginInput) (string, *model.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, &AppError{Code: http.StatusUnauthorized, Message: "Email atau password salah"}
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", nil, &AppError{Code: http.StatusUnauthorized, Message: "Email atau password salah"}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Duration(s.cfg.JWTExpireHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, user, nil
}

func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &AppError{Code: http.StatusNotFound, Message: "Pengguna tidak ditemukan"}
		}
		return nil, err
	}
	return user, nil
}
