package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	"siminventaris-backend/internal/middleware"
)

// newBarangTestApp memasang route barang persis seperti router.RegisterBarangRoutes
// tanpa memerlukan koneksi database (service nil karena handler tidak sampai ke situ).
func newBarangTestApp(t *testing.T, secret string) *fiber.App {
	t.Helper()
	app := fiber.New()
	h := NewBarangHandler(nil)
	barang := app.Group("/api/barang", middleware.WajibLogin(secret))
	barang.Post("/", middleware.HanyaPeran("admin"), h.Create)
	barang.Put("/:id", middleware.HanyaPeran("admin"), h.Update)
	barang.Delete("/:id", middleware.HanyaPeran("admin"), h.Delete)
	return app
}

func TestBarangHandlerWriteRequiresLoginAndAdmin(t *testing.T) {
	const secret = "rahasia_uji_cukup_panjang_1234567890"
	app := newBarangTestApp(t, secret)
	validBody := `{"kode_barang":"LAB-001","nama":"Oscilloscope","stok":1,"kondisi":"baik","kategori_id":1}`

	// Tanpa token -> 401
	req := httptest.NewRequest(http.MethodPost, "/api/barang", strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", resp.StatusCode)
	}

	// Token member -> 403 (hanya admin boleh tulis)
	req = httptest.NewRequest(http.MethodPost, "/api/barang", strings.NewReader(validBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, secret, 2, "member"))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403 for member, got %d", resp.StatusCode)
	}
}

func TestBarangHandlerCreateBodyAndValidation(t *testing.T) {
	const secret = "rahasia_uji_cukup_panjang_1234567890"
	app := newBarangTestApp(t, secret)

	// Body kosong dengan token admin -> 400 (bukan 422)
	req := httptest.NewRequest(http.MethodPost, "/api/barang", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, secret, 1, "admin"))
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

	// JSON valid tetapi field melanggar aturan validasi -> 422
	req = httptest.NewRequest(http.MethodPost, "/api/barang",
		strings.NewReader(`{"kode_barang":"A","nama":"","stok":-1,"kondisi":"jelek","kategori_id":0}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, secret, 1, "admin"))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for invalid fields, got %d", resp.StatusCode)
	}
	out := decodeResponse(t, resp)
	if out.Sukses {
		t.Fatal("expected sukses=false on validation error")
	}
	if out.Data == nil {
		t.Fatal("expected validation detail map on 422")
	}
}

func TestBarangHandlerUpdateDeleteRejectInvalidID(t *testing.T) {
	const secret = "rahasia_uji_cukup_panjang_1234567890"
	app := newBarangTestApp(t, secret)
	auth := "Bearer " + testToken(t, secret, 1, "admin")

	req := httptest.NewRequest(http.MethodPut, "/api/barang/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for non-numeric update id, got %d", resp.StatusCode)
	}
	if out := decodeResponse(t, resp); out.Pesan != "ID barang tidak valid" {
		t.Fatalf("unexpected message: %q", out.Pesan)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/barang/abc", nil)
	req.Header.Set("Authorization", auth)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for non-numeric delete id, got %d", resp.StatusCode)
	}
}
