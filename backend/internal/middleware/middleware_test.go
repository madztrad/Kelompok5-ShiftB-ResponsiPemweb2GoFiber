package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func createTestToken(secret string, userID uint, role string, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     exp.Unix(),
	})
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestWajibLoginAndHanyaPeran(t *testing.T) {
	const secret = "test_secret_yang_cukup_panjang"

	app := fiber.New()
	app.Get("/protected", WajibLogin(secret), func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	app.Delete("/admin-only", WajibLogin(secret), HanyaPeran("admin"), func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	// 1. Tanpa token -> 401
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	// 2. Token kedaluwarsa -> 401
	expired := createTestToken(secret, 1, "member", time.Now().Add(-1*time.Hour))
	reqExpired := httptest.NewRequest(http.MethodGet, "/protected", nil)
	reqExpired.Header.Set("Authorization", "Bearer "+expired)
	respExpired, _ := app.Test(reqExpired)
	if respExpired.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", respExpired.StatusCode)
	}

	// 3. Member akses route admin -> 403
	memberToken := createTestToken(secret, 2, "member", time.Now().Add(time.Hour))
	reqMember := httptest.NewRequest(http.MethodDelete, "/admin-only", nil)
	reqMember.Header.Set("Authorization", "Bearer "+memberToken)
	respMember, _ := app.Test(reqMember)
	if respMember.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for member, got %d", respMember.StatusCode)
	}

	// 4. Admin akses route admin -> 200
	adminToken := createTestToken(secret, 1, "admin", time.Now().Add(time.Hour))
	reqAdmin := httptest.NewRequest(http.MethodDelete, "/admin-only", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	respAdmin, _ := app.Test(reqAdmin)
	if respAdmin.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for admin, got %d", respAdmin.StatusCode)
	}
}
