package utils

import "testing"

type sampleInput struct {
	Nama  string `validate:"required,min=2"`
	Email string `validate:"required,email"`
	Stok  int    `validate:"gte=0"`
}

func TestValidateStruct(t *testing.T) {
	invalid := sampleInput{
		Nama:  "",
		Email: "bukan-email",
		Stok:  -1,
	}
	errs := ValidateStruct(invalid)
	if errs == nil {
		t.Fatal("expected validation errors, got nil")
	}
	if errs["nama"] == "" || errs["email"] == "" || errs["stok"] == "" {
		t.Fatalf("expected errors for nama, email, and stok, got: %v", errs)
	}

	valid := sampleInput{
		Nama:  "Router",
		Email: "admin@lab.test",
		Stok:  5,
	}
	if errs := ValidateStruct(valid); errs != nil {
		t.Fatalf("expected nil errors, got: %v", errs)
	}
}
