package apiv1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (api *ApiHandler) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing form during register")
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

	err = api.databaseConnection.RegisterNewUser(ctx, username, password)
	if err == database.ErrOnCreateUserExists {
		log.Info().Str("Username", username).Msg("User already exists")
		w.Write([]byte("Username already exists!"))
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Error during register of new user")
		w.Write([]byte("An error occurred during registration of user, please try again later."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info().Err(err).Msg("User registered")
	w.Write([]byte("Registration successful!"))
	w.WriteHeader(http.StatusCreated)
}
