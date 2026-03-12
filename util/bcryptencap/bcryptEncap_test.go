package bcryptencap

import (
	"testing"
)

func TestGenerateFromPassword(t *testing.T) {
	password := "testPassword"
	hash, err := GenerateFromPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate hash from password: %v", err)
	}

	if len(hash) == 0 {
		t.Fatal("Expected hash to be not empty")
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	password := "testPassword"
	hash, _ := GenerateFromPassword(password)

	err := CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Fatalf("Failed to compare hash and password: %v", err)
	}
}
