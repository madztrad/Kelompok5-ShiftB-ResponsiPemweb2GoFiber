package service

import (
	"errors"
	"net/http"
	"testing"

	"siminventaris-backend/internal/utils"
)

func TestValidatePeminjamanUpdateAccess(t *testing.T) {
	jumlah := 2
	tanggalKembali := "2026-11-28"
	catatan := "Untuk praktikum"

	tests := []struct {
		name    string
		role    string
		input   UpdatePeminjamanInput
		wantErr bool
	}{
		{name: "member can edit own note", role: "member", input: UpdatePeminjamanInput{Catatan: &catatan}},
		{name: "member cannot change quantity", role: "member", input: UpdatePeminjamanInput{Jumlah: &jumlah}, wantErr: true},
		{name: "member cannot change dates", role: "member", input: UpdatePeminjamanInput{TanggalKembali: tanggalKembali}, wantErr: true},
		{name: "member cannot change status", role: "member", input: UpdatePeminjamanInput{Status: "dikembalikan"}, wantErr: true},
		{name: "admin can edit all fields", role: "admin", input: UpdatePeminjamanInput{Jumlah: &jumlah, TanggalKembali: tanggalKembali, Status: "dikembalikan"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePeminjamanUpdateAccess(tt.role, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validatePeminjamanUpdateAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				var appErr *AppError
				if !errors.As(err, &appErr) || appErr.Code != http.StatusForbidden {
					t.Fatalf("expected 403 AppError, got %v", err)
				}
			}
		})
	}
}

func TestPeminjamanDatesMustUseISODateFormat(t *testing.T) {
	valid := CreatePeminjamanInput{
		BarangID:      1,
		Jumlah:        1,
		TanggalPinjam: "2026-11-27",
	}
	if errs := utils.ValidateStruct(valid); errs != nil {
		t.Fatalf("expected valid date, got validation errors: %v", errs)
	}

	invalid := valid
	invalid.TanggalPinjam = "besok"
	if errs := utils.ValidateStruct(invalid); errs == nil {
		t.Fatal("expected invalid date validation error")
	}

	emptyNote := ""
	update := UpdatePeminjamanInput{
		Catatan: &emptyNote,
	}
	if errs := utils.ValidateStruct(update); errs != nil {
		t.Fatalf("expected empty optional note to be valid, got validation errors: %v", errs)
	}

	update.TanggalKembali = "not-a-date"
	if errs := utils.ValidateStruct(update); errs == nil {
		t.Fatal("expected invalid return date validation error")
	}
}
