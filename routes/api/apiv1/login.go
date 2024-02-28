package apiv1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (api *ApiHandler) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing form during login")
		w.Write([]byte("Could not parse form!"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.Form.Get("username")
	if username == "" {
		log.Info().Msg("Request did not include 'username' field!")
		w.Write([]byte("Request must include 'username' field!"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password := r.Form.Get("password")
	if password == "" {
		log.Info().Msg("Request did not include 'password' field!")
		w.Write([]byte("Request must include 'password' field!"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(password) > passwordMaxLen {
		log.Info().Int("PasswordLen", len(password)).Msg("Password is too long! ")
		w.Write([]byte(fmt.Sprintf("Password must be less than %v characters long!", passwordMaxLen)))
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	ctx := context.Background()
	ok, err := api.databaseConnection.ValidateAuthenticationAttempt(ctx, username, password)
	if err == database.ErrOnFetchUserDoesNotExist {
		log.Info().Str("Username", username).Msg("Login Request for Invalid Username")
		w.Write([]byte("No user with given username exists!"))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Found error during authentication!")
		w.Write([]byte("An error occurred during authentication attempt, please try again."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Actually check if the user is who they say they are
	if !ok {
		log.Info().Str("Username", username).Msg("Invalid login attempt!")
		w.Write([]byte("Invalid Username or Password."))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Now we can go about giving the JWT to authenticate in the future

	token, err := api.generateJWT(username)
	if err != nil {
		log.Error().Err(err).Str("Username", username).Msg("Error during creation of JWT!")
		w.Write([]byte("An error occurred during authentication attempt, please try again"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
	w.WriteHeader(http.StatusOK)
}
