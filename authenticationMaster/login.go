package authenticationmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
	"github.com/rs/zerolog/log"
)

func (authMaster *AuthenticationMaster) Login(w http.ResponseWriter, r *http.Request) {
	var requestCredentials httpRequestCredentials
	err := json.NewDecoder(r.Body).Decode(&requestCredentials)
	if err != nil {
		log.Error().Err(err).Msg("Found error parsing request during login")
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

	ok, err := authMaster.databaseConnection.ValidateLoginAttempt(databaseQueryContext, requestCredentials.Username, requestCredentials.Password)
	if databaseQueryContext.Err() == context.DeadlineExceeded {
		log.Info().Str("Username", requestCredentials.Username).Msg("Database query duration exceeded!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error!"))
		return
	}
	if err == database.ErrOnFetchUserDoesNotExist {
		log.Info().Str("Username", requestCredentials.Username).Msg("Login Request for Invalid Username")
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
		log.Info().Str("Username", requestCredentials.Username).Msg("Invalid login attempt!")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	// Now we can go about giving the JWT to authenticate in the future

	userID, _ := authMaster.databaseConnection.GetUserIDByUsername(context.Background(), requestCredentials.Username)
	token, err := authMaster.generateJWT(userID)
	if err != nil {
		log.Error().Err(err).Str("Username", requestCredentials.Username).Msg("Error during creation of JWT!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication attempt, please try again"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
