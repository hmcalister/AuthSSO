package authenticationmaster

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	issuerString       string            = "hmcalisterAuthSSO"
	tokenSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256
)

// Given a UserID, generate a new token with that userID as the subject
func (authMaster *AuthenticationMaster) generateJWT(userID string) (string, error) {
	currentTime := time.Now()
	expirationTime := currentTime.Add(tokenExpirationDuration)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuerString,
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userID,
		},
	)
	signedString, err := token.SignedString(authMaster.tokenSecretKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
