package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"siminventaris-backend/internal/middleware"
	"siminventaris-backend/internal/utils"
)

func testToken(t *testing.T, secret string, userID uint, role string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("gagal membuat token: %v", err)
	}
	return signed
}

func decodeResponse(t *testing.T, resp *http.Response) utils.APIResponse {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("gagal membaca body: %v", err)
	}
	var out utils.APIResponse
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("body bukan JSON valid: %s", body)
	}
	return out
}

func TestAuthHandlerRegisterBadBody(t *testing.T) {
	app := fiber.New()
	app.Post("/api/auth/register", NewAuthHandler(nil).Register)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for empty body, got %d", resp.StatusCode)
	}
	if out := decodeResponse(t, resp); out.Pesan != "Format request body tidak valid" {
		t.Fatalf("unexpected message: %q", out.Pesan)
	}
}

func TestAuthHandlerProfileWithoutToken(t *testing.T) {
	const secret = "rahasia_uji_cukup_panjang_1234567890"
	app := fiber.New()
	app.Get("/api/auth/profile", middleware.WajibLogin(secret), NewAuthHandler(nil).Profile)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", resp.StatusCode)
	}
}

func TestKategoriHandlerValidationAndAuth(t *testing.T) {
	const secret = "rahasia_uji_cukup_panjang_1234567890"
	app := fiber.New()
	app.Post("/api/kategori", middleware.WajibLogin(secret), middleware.HanyaPeran("admin"), NewKategoriHandler(nil).Create)

	req := httptest.NewRequest(http.MethodPost, "/api/kategori", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/kategori", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, secret, 2, "member"))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403 for member, got %d", resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/kategori", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, secret, 1, "admin"))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for empty body, got %d", resp.StatusCode)
	}
}
