package authenticationmaster

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type authorizedUserData struct {
	UserID   string
	Username string
}

// Authenticate a request by checking the JWT in the request header (or cookie).
//
// Note this effectively reimplements the logic of the go-chi jwtauth Verifier middleware,
// but exposes the logic on a route rather than as middleware. https://pkg.go.dev/github.com/go-chi/jwtauth/v5@v5.3.0#Verifier
func (authMaster *AuthenticationMaster) AuthenticateRequest(w http.ResponseWriter, r *http.Request) {
	token, err := jwtauth.VerifyRequest(authMaster.tokenAuth, r, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)

	if token == nil {
		slog.Debug("No token received.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// The possible error return types are actually fairly limited,
	// based on the source code of VerifyToken and ErrorReason
	//
	// https://pkg.go.dev/github.com/go-chi/jwtauth/v5@v5.3.0#VerifyToken
	// https://pkg.go.dev/github.com/go-chi/jwtauth/v5@v5.3.0#ErrorReason
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
		slog.Error("UserID does not exist in database", "Error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Send the result as the response
	userData := authorizedUserData{
		UserID:   userID,
		Username: username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
