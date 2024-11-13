package authenticationmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
)

func (authMaster *AuthenticationMaster) Register(w http.ResponseWriter, r *http.Request) {
	var requestCredentials httpRequestCredentials
	err := json.NewDecoder(r.Body).Decode(&requestCredentials)
	if err != nil {
		slog.Error("Found error parsing request during register", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse form!"))
		return
	}

	if requestCredentials.Username == "" {
		slog.Info("Request did not include 'username' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'username' field!"))
		return
	}

	sanitizedUsername := authMaster.htmlSanitizer.Sanitize(requestCredentials.Username)
	if requestCredentials.Username != sanitizedUsername {
		slog.Info("Username must not require sanitization!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username must not require sanitization!"))
		return
	}

	if requestCredentials.Password == "" {
		slog.Info("Request did not include 'password' field!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must include 'password' field!"))
		return
	}
	if len(requestCredentials.Password) > passwordMaxLen {
		slog.Info("Password is too long!", "PasswordLength", len(requestCredentials.Password))
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(fmt.Sprintf("Password must be less than %v characters long!", passwordMaxLen)))
		return
	}

	databaseQueryContext, databaseQueryContextCancel := context.WithTimeout(context.Background(), maximumDatabaseQueryDuration)
	defer databaseQueryContextCancel()

	err = authMaster.databaseConnection.RegisterNewUser(databaseQueryContext, requestCredentials.Username, requestCredentials.Password)
	if databaseQueryContext.Err() == context.DeadlineExceeded {
		slog.Info("Database query duration exceeded!", "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error!"))
		return
	}
	if err == database.ErrOnCreateUserExists {
		slog.Info("User already exists", "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username already exists!"))
		return
	}
	if err != nil {
		slog.Error("Error during register of new user", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during registration of user, please try again later."))
		return
	}

	slog.Info("User registered", "Username", requestCredentials.Username)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration successful!"))
}
