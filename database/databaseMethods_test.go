package database_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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
	os.Remove(testDatabasePath)
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

func TestSequentialDatabaseAccess(t *testing.T) {
	numUsers := 256

	password := "Password123"
	ctx := context.Background()
	for i := 0; i < numUsers; i += 1 {
		username := fmt.Sprintf("User%v", i)

		err := databaseManager.RegisterNewUser(ctx, username, password)
		if err != nil {
			t.Errorf("Failed to register user %v: %v", i, err)
		}
	}
	for i := 0; i < numUsers; i += 1 {
		username := fmt.Sprintf("User%v", i)
		ok, err := databaseManager.ValidateAuthenticationAttempt(ctx, username, password)
		if err != nil {
			t.Errorf("Error during authentication of user %v: %v", i, err)
		}
		if !ok {
			t.Errorf("Failed to authenticate user %v (with correct credentials)", i)
		}
	}
	for i := 0; i < numUsers; i += 1 {
		username := fmt.Sprintf("User%v", i)
		err := databaseManager.DeleteUserByUsername(ctx, username)
		if err != nil {
			t.Errorf("Failed to delete user %v: %v", i, err)
		}
	}
}

func TestParallelDatabaseAccess(t *testing.T) {
	numUsers := 1024
	numWorkers := 10

	var wg sync.WaitGroup
	workerContext, workerCancel := context.WithCancel(context.Background())
	errorChan := make(chan error)
	userChan := make(chan int)

	for i := 0; i < numWorkers; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx := context.Background()
			for {
				select {
				case <-workerContext.Done():
					return
				case i := <-userChan:
					username := fmt.Sprintf("user%v", i)
					err := databaseManager.RegisterNewUser(ctx, username, "Password123")
					if err != nil {
						errorChan <- err
						workerCancel()
					}

					_, err = databaseManager.ValidateAuthenticationAttempt(ctx, username, "Password123")
					if err != nil {
						errorChan <- err
						workerCancel()
					}

					// err = databaseManager.DeleteUserByUsername(ctx, username)
					// if err != nil {
					// 	errorChan <- err
					// 	workerCancel()
					// }
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	go func() {
		for i := 0; i < numUsers; i += 1 {
			userChan <- i
		}
		workerCancel()
	}()

	for err := range errorChan {
		t.Errorf("Failed to interact with database in a parallel fashion: %v", err)
	}
}
