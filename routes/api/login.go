package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (api *ApiHandler) login(w http.ResponseWriter, r *http.Request) {
	var user userJsonData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing request during login")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse form!"))
		return
	}

	if user.Username == "" {
		log.Info().Msg("Request did not include 'username' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'username' field!"))
		return
	}

	if user.Password == "" {
		log.Info().Msg("Request did not include 'password' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'password' field!"))
		return
	}
	if len(user.Password) > passwordMaxLen {
		log.Info().Int("PasswordLen", len(user.Password)).Msg("Password is too long! ")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(fmt.Sprintf("Password must be less than %v characters long!", passwordMaxLen)))
		return
	}

	ctx := context.Background()
	ok, err := api.databaseConnection.ValidateAuthenticationAttempt(ctx, user.Username, user.Password)
	if err == database.ErrOnFetchUserDoesNotExist {
		log.Info().Str("Username", user.Username).Msg("Login Request for Invalid Username")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No user with given username exists!"))
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Found error during authentication!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication attempt, please try again."))
		return
	}

	// Actually check if the user is who they say they are
	if !ok {
		log.Info().Str("Username", user.Username).Msg("Invalid login attempt!")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	// Now we can go about giving the JWT to authenticate in the future

	token, err := api.generateJWT(user.Username)
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("Error during creation of JWT!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication attempt, please try again"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
