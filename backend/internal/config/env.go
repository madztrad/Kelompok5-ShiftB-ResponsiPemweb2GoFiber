package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	minJWTSecretLength = 32
	legacyDevSecret    = "dev_only_secret_ganti_di_env"
)

type EnvConfig struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	CORSOrigins    []string
}

func LoadEnv() EnvConfig {
	// Abaikan error jika file .env tidak ada (misal di platform hosting)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: file .env tidak ditemukan, membaca environment variable sistem")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "3000"
	}

	expireHours, err := strconv.Atoi(getEnvOrDefault("JWT_EXPIRE_HOURS", "24"))
	if err != nil || expireHours <= 0 {
		expireHours = 24
	}

	rawOrigins := getEnvOrDefault("CORS_ORIGINS", "http://localhost:5500,http://127.0.0.1:5500")
	var origins []string
	for _, origin := range strings.Split(rawOrigins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < minJWTSecretLength || jwtSecret == legacyDevSecret {
		log.Fatalf("JWT_SECRET wajib diisi dengan string acak minimal %d karakter (jangan memakai nilai bawaan)", minJWTSecretLength)
	}

	return EnvConfig{
		Port:           port,
		DBHost:         getEnvOrDefault("DB_HOST", "127.0.0.1"),
		DBPort:         getEnvOrDefault("DB_PORT", "3306"),
		DBUser:         getEnvOrDefault("DB_USER", "root"),
		DBPass:         os.Getenv("DB_PASS"),
		DBName:         getEnvOrDefault("DB_NAME", "siminventaris_db"),
		JWTSecret:      jwtSecret,
		JWTExpireHours: expireHours,
		CORSOrigins:    origins,
	}
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
