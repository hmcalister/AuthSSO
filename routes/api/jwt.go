package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpirationDuration time.Duration = 6 * time.Hour
)

var (
	issuerString       string            = "hmcalisterAuthSSO"
	tokenSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256
)

func (api *ApiHandler) generateJWT(username string) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(tokenExpirationDuration)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuerString,
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   username,
		},
	)
	signedString, err := token.SignedString(api.tokenSecretKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
