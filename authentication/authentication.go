package authentication

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/hmcalister/AuthSSO/database/database"
	"golang.org/x/crypto/argon2"
)

// Defines the parameters to argon2.IDKey to be used for authentication
type AuthenticationParams struct {
	saltLen    uint32
	timeCost   uint32
	memoryCost uint32
	threads    uint8
	keyLen     uint32
}

// Generate a new salt and return it.
//
// This function can error if a kernel function errors, although this should never happen.
func (authParams *AuthenticationParams) generateSalt() (string, error) {
	salt := make([]byte, authParams.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	return string(salt), nil
}

// Perform the hash of a (plaintext) password with salt.
func (authParams *AuthenticationParams) calculateHash(password string, salt string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), authParams.timeCost, authParams.memoryCost, authParams.threads, authParams.keyLen)
	return string(hash)
}

// Given the newly generated UUID and the incoming (plaintext) password, generate a new AuthenticationDatum
//
// Note this function generates a new salt, so should only be called once during registration.
// Also note this function does NOT save the authenticationdatum to the database, this should be done after this function call.
func (authParams *AuthenticationParams) NewAuthenticationDatum(uuid string, password string) (*database.AuthenticationDatum, error) {
	salt, err := authParams.generateSalt()
	if err != nil {
		return nil, err
	}

	hashedPassword := authParams.calculateHash(password, salt)

	return &database.AuthenticationDatum{
		Uuid:           uuid,
		Hashedpassword: string(hashedPassword),
		Salt:           string(salt),
	}, nil
}

// Given an AuthenticationDatum and a (plaintext) password attempt, check to see if the password matches the hash as requested
//
// Returns true if the password matches the hash, false otherwise.
// This function will error if len(salt) is not equal to auth_SALT_LEN constant, or if len(authDatum.Hashedpassword) is not equal to authParams.keyLen
func (authParams *AuthenticationParams) ValidateAuthenticationAttempt(authDatum *database.AuthenticationDatum, passwordAttempt string) (bool, error) {
	if len(authDatum.Salt) != int(authParams.saltLen) {
		return false, errors.New("length of authDatum salt does not equal authParams saltLen")
	}

	if len(authDatum.Hashedpassword) != int(authParams.keyLen) {
		return false, errors.New("length of authDatum hashedPassword does not equal authParams keyLen")
	}

	attemptHash := authParams.calculateHash(passwordAttempt, authDatum.Salt)
	return authDatum.Hashedpassword == attemptHash, nil
}
