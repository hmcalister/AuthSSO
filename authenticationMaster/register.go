package authenticationmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (authMaster *AuthenticationMaster) Register(w http.ResponseWriter, r *http.Request) {
	var requestCredentials httpRequestCredentials
	err := json.NewDecoder(r.Body).Decode(&requestCredentials)
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing request during register")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse form!"))
		return
	}

	if requestCredentials.Username == "" {
		log.Info().Msg("Request did not include 'username' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'username' field!"))
		return
	}

	sanitizedUsername := authMaster.htmlSanitizer.Sanitize(requestCredentials.Username)
	if requestCredentials.Username != sanitizedUsername {
		log.Info().Msg("Username must not require sanitization!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username must not require sanitization!"))
		return
	}

	if requestCredentials.Password == "" {
		log.Info().Msg("Request did not include 'password' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'password' field!"))
		return
	}
	if len(requestCredentials.Password) > passwordMaxLen {
		log.Info().Int("PasswordLen", len(requestCredentials.Password)).Msg("Password is too long! ")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(fmt.Sprintf("Password must be less than %v characters long!", passwordMaxLen)))
		return
	}

	databaseQueryContext, databaseQueryContextCancel := context.WithTimeout(context.Background(), maximumDatabaseQueryDuration)
	defer databaseQueryContextCancel()

	err = authMaster.databaseConnection.RegisterNewUser(databaseQueryContext, requestCredentials.Username, requestCredentials.Password)
	if databaseQueryContext.Err() == context.DeadlineExceeded {
		log.Info().Str("Username", requestCredentials.Username).Msg("Database query duration exceeded!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error!"))
		return
	}
	if err == database.ErrOnCreateUserExists {
		log.Info().Str("Username", requestCredentials.Username).Msg("User already exists")
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
