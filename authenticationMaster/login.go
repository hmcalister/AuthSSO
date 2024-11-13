package authenticationmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hmcalister/AuthSSO/database"
)

func (authMaster *AuthenticationMaster) Login(w http.ResponseWriter, r *http.Request) {
	var requestCredentials httpRequestCredentials
	err := json.NewDecoder(r.Body).Decode(&requestCredentials)
	if err != nil {
		slog.Error("Found error parsing request during login", "Error", err)
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

	ok, err := authMaster.databaseConnection.ValidateLoginAttempt(databaseQueryContext, requestCredentials.Username, requestCredentials.Password)
	if databaseQueryContext.Err() == context.DeadlineExceeded {
		slog.Info("Database query duration exceeded!", "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error!"))
		return
	}
	if err == database.ErrOnFetchUserDoesNotExist {
		slog.Info("Login Request for Invalid Username", "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No user with given username exists!"))
		return
	}
	if err != nil {
		slog.Error("Found error during authentication!", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication attempt, please try again."))
		return
	}

	// Actually check if the user is who they say they are
	if !ok {
		slog.Info("Invalid login attempt!", "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	// Now we can go about giving the JWT to authenticate in the future

	userID, _ := authMaster.databaseConnection.GetUserIDByUsername(context.Background(), requestCredentials.Username)
	token, err := authMaster.generateJWT(userID)
	if err != nil {
		slog.Error("Error during creation of JWT!", "Error", err, "Username", requestCredentials.Username)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred during authentication attempt, please try again"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
