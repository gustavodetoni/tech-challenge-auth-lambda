package document_test

import (
	"testing"

	"github.com/soat-architecture/tech-challenge-auth-lambda/pkg/br/document"
)

func TestNormalize(t *testing.T) {
	got := document.Normalize("123.456.789-09")
	if got != "12345678909" {
		t.Fatalf("Normalize() = %q", got)
	}
}

func TestIsValidCPF(t *testing.T) {
	if !document.IsValidCPF("529.982.247-25") {
		t.Fatal("expected valid CPF")
	}
	if document.IsValidCPF("111.111.111-11") {
		t.Fatal("expected repeated CPF to be invalid")
	}
}

func TestIsValidCNPJ(t *testing.T) {
	if !document.IsValidCNPJ("04.252.011/0001-10") {
		t.Fatal("expected valid CNPJ")
	}
	if document.IsValidCNPJ("00.000.000/0000-00") {
		t.Fatal("expected repeated CNPJ to be invalid")
	}
}
