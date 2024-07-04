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

}
