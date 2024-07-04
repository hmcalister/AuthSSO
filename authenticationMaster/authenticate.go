package authenticationmaster

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"
)

type authorizedUserData struct {
	UserID   string
	Username string
}

func (authMaster *AuthenticationMaster) AuthenticateRequest(w http.ResponseWriter, r *http.Request) {
	token, err := jwtauth.VerifyRequest(authMaster.tokenAuth, r, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)

	if token == nil {
		log.Debug().Msg("No token received.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		switch err {
		case jwtauth.ErrExpired:
			w.Write([]byte("Token is expired."))
		case jwtauth.ErrIATInvalid:
			w.Write([]byte("Token issued time invalid."))
		case jwtauth.ErrNBFInvalid:
			w.Write([]byte("Token not yet valid."))
		default:
			w.Write([]byte("Token unauthorized."))
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Extract the UserID from the token
	userID := token.Subject()

	// Query the database and get the username from it
	ctx := context.Background()
	username, err := authMaster.databaseConnection.GetUsernameByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("UserID does not exist in database.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}
