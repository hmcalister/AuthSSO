package database_test

import (
	"context"
	"log"
	"testing"

	"github.com/hmcalister/AuthSSO/database"
)

const (
	testDatabasePath string = "databaseTest.sqlite"
)

var (
	databaseManager *database.DatabaseManager
)

func TestMain(m *testing.M) {
	var err error
	databaseManager, err = database.NewDatabase(testDatabasePath)

	ctx := context.Background()
	databaseManager.RegisterNewUser(ctx, "John Smith", "Password123")
	databaseManager.RegisterNewUser(ctx, "Jane Doe", "QWERTY321")

	if err != nil {
		log.Fatalf("encountered error when opening test database %v", err)
	}

	m.Run()

	err = databaseManager.CloseDatabase()
	if err != nil {
		log.Fatalf("encountered error while closing database %v", err)
	}
	// os.Remove(testDatabasePath)
}

func TestRegisterNewUser(t *testing.T) {
	ctx := context.Background()
	err := databaseManager.RegisterNewUser(ctx, "newUser", "Password123")
	if err != nil {
		t.Errorf("Error while registering new user: %v", err)
	}
}

func TestAttemptRegisterExistingUser(t *testing.T) {
	ctx := context.Background()
	err := databaseManager.RegisterNewUser(ctx, "John Smith", "Password123")
	if err == nil {
		t.Errorf("No error thrown while registering a new user with exact credentials: %v", err)
	}

	err = databaseManager.RegisterNewUser(ctx, "John Smith", "newPassword")
	if err == nil {
		t.Errorf("No error thrown while registering a new user with same username, different password: %v", err)
	}
}

func TestDeleteExistingUser(t *testing.T) {
	ctx := context.Background()
	err := databaseManager.DeleteUserByUsername(ctx, "Jane Doe")
	if err != nil {
		t.Fatalf("Error while deleting user: %v", err)
	}

	err = databaseManager.RegisterNewUser(ctx, "Jane Doe", "QWERTY321")
	if err != nil {
		t.Errorf("Error while recreating deleted user (checking user is *actually* deleted): %v", err)
	}
}

func TestValidateAuthenticationAttempt(t *testing.T) {
	ctx := context.Background()

	valid, err := databaseManager.ValidateAuthenticationAttempt(ctx, "John Smith", "Password123")
	if err != nil {
		t.Errorf("Error while authenticating (with correct credentials): %v", err)
	}
	if valid == false {
		t.Errorf("Authentication attempt failed (when presented with correct credentials)")
	}

	valid, err = databaseManager.ValidateAuthenticationAttempt(ctx, "John Smith", "IncorrectPassword")
	if err != nil {
		t.Errorf("Error while authenticating (with incorrect credentials): %v", err)
	}
	if valid == true {
		t.Errorf("Authentication attempt succeeded (when presented with incorrect credentials)")
	}
}
