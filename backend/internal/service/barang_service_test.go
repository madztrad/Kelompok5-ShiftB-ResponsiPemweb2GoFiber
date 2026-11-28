package service

import (
	"testing"

	"siminventaris-backend/internal/utils"
)

func TestBarangInputValidation(t *testing.T) {
	valid := BarangInput{
		KodeBarang: "LAB-001",
		Nama:       "Oscilloscope Digital",
		Stok:       3,
		Kondisi:    "baik",
		KategoriID: 1,
	}
	if errs := utils.ValidateStruct(valid); errs != nil {
		t.Fatalf("expected valid input, got errors: %v", errs)
	}

	tests := []struct {
		name  string
		mutat func(*BarangInput)
		field string
	}{
		{name: "kode terlalu pendek", mutat: func(b *BarangInput) { b.KodeBarang = "A" }, field: "kode_barang"},
		{name: "nama kosong", mutat: func(b *BarangInput) { b.Nama = "" }, field: "nama"},
		{name: "stok negatif", mutat: func(b *BarangInput) { b.Stok = -1 }, field: "stok"},
		{name: "kondisi tidak valid", mutat: func(b *BarangInput) { b.Kondisi = "jelek" }, field: "kondisi"},
		{name: "kategori nol", mutat: func(b *BarangInput) { b.KategoriID = 0 }, field: "kategori_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.mutat(&input)
			errs := utils.ValidateStruct(input)
			if errs == nil {
				t.Fatalf("expected validation error for %s", tt.field)
			}
			if _, ok := errs[tt.field]; !ok {
				t.Fatalf("expected error on field %s, got %v", tt.field, errs)
			}
		})
	}
}
