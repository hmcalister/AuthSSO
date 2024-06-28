package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (api *ApiHandler) register(w http.ResponseWriter, r *http.Request) {
	var user userJsonData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing request during register")
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

	err = api.databaseConnection.RegisterNewUser(ctx, user.Username, user.Password)
	if err == database.ErrOnCreateUserExists {
		log.Info().Str("Username", user.Username).Msg("User already exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username already exists!"))
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Error during register of new user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during registration of user, please try again later."))
		return
	}

	log.Info().Err(err).Msg("User registered")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration successful!"))
}
