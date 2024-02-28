package apiv1

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	expirationDuration time.Duration = 1 * time.Hour
)

var (
	tokenSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256
)

func (api *ApiHandler) generateJWT(username string) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(expirationDuration)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
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
