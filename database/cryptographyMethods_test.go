package database

import (
	"testing"
)

func TestSaltGeneration(t *testing.T) {
	_, err := generateSalt()
	if err != nil {
		t.Errorf("Error during salt generation %v", err)
	}
}

func TestHashCalculation(t *testing.T) {

	password1 := "password123"
	password2 := "qwerty321"
	salt1, _ := generateSalt()
	salt2, _ := generateSalt()

	if calculateHash(password1, salt1) == calculateHash(password1, salt2) {
		t.Error("Hashes of same password with different salt are equal")
	}

	if calculateHash(password1, salt1) == calculateHash(password2, salt1) {
		t.Error("Hashes of different password with same salt are equal")
	}

	if calculateHash(password1, salt1) != calculateHash(password1, salt1) {
		t.Error("Hashes of same password and same salt are not equal")
	}
}
