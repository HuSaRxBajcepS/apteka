package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testp"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() zwróciło błąd: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() zwróciło pusty hash")
	}

	if hash == password {
		t.Fatal("Hash nie powinien być równy oryginalnemu hasłu")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Fatal("Wygenerowany hash nie odpowiada podanemu hasłu")
	}
	t.Log("Hash:", hash)
}
