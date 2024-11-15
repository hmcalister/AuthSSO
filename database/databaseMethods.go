package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hmcalister/AuthSSO/database/sqlc"

	_ "embed"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

type DatabaseManager struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// Creates a new database struct at the given filepath (erroring if not possible)
// and creating the schema detailed by schema.sql and query.sql (if not already present)
func NewDatabase(databaseFilePath string) (*DatabaseManager, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", databaseFilePath)
	if err != nil {
		return nil, err
	}

	// Massively improves parallelism
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	return &DatabaseManager{
		db:      db,
		queries: queries,
	}, nil
}

func (database *DatabaseManager) CloseDatabase() error {
	return database.db.Close()
}

// Checks if a user exists in the database. Returns true if the user exists already.
func (database *DatabaseManager) CheckUserExists(ctx context.Context, username string) (bool, error) {
	user, err := database.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return user.Username == "", nil
}

// Gets the username of the user associated with a specific userID. Returns an error if the userID does not exist.
func (database *DatabaseManager) GetUsernameByUserID(ctx context.Context, userID string) (string, error) {
	user, err := database.queries.GetUserByUUID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

// Gets the userID of a user specified by the username. Returns an error if the username does not exist.
func (database *DatabaseManager) GetUserIDByUsername(ctx context.Context, username string) (string, error) {
	user, err := database.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	return user.Uuid, nil
}

// Given a username and password, attempt to register a new user
//
// Fails (and returns a non-nil error) if:
// - The salt fails to be generated
// - The username already exists in the database (ErrUserExists)
// - The transaction to store both the new user data and the new auth data fails
//
// This method ensures that the new user data and auth data is create atomically, so
// a user cannot exist without auth data, and auth data cannot exist without a user
func (database *DatabaseManager) RegisterNewUser(ctx context.Context, username string, password string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	hashedPassword := calculateHash(password, salt)
	newUserUUID := uuid.New().String()

	newUserAuthDatum := sqlc.CreateAuthenticationDataParams{
		Uuid:           newUserUUID,
		HashedPassword: string(hashedPassword),
		Salt:           string(salt),
	}

	newUserDatum := sqlc.CreateUserParams{
		Uuid:     newUserUUID,
		Username: username,
	}

	// Begin database transaction to ensure user and authdata created together
	tx, err := database.db.Begin()
	if err != nil {
		return err
	}
	// If anything fails, an early return is triggered (before tx.Commit is called) and tx.Rollback is called
	defer tx.Rollback()

	qtx := database.queries.WithTx(tx)
	_, err = qtx.CreateAuthenticationData(ctx, newUserAuthDatum)
	if err != nil {
		return err
	}
	_, err = qtx.CreateUser(ctx, newUserDatum)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch errCode := sqliteErr.ExtendedCode; errCode {
		case sqlite3.ErrConstraintUnique:
			return ErrOnCreateUserExists
		}
	}

	if err != nil {
		return err
	}
	return tx.Commit()
}

// Delete a user from the database, including the authdata and user.
//
// Fails and returns a non-nil error if:
// - The user does not exist in the database
// - The transaction to delete both the user data and auth data fails
func (database *DatabaseManager) DeleteUserByUsername(ctx context.Context, username string) error {
	// Get the user by username, if it exists
	userData, err := database.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrOnFetchUserDoesNotExist
		}
		return err
	}

	userUUID := userData.Uuid

	// Begin database transaction to ensure user and authdata deleted together
	tx, err := database.db.Begin()
	if err != nil {
		return err
	}
	// If anything fails, an early return is triggered (before tx.Commit is called) and tx.Rollback is called
	defer tx.Rollback()

	qtx := database.queries.WithTx(tx)
	err = qtx.DeleteUser(ctx, userUUID)
	if err != nil {
		return err
	}
	err = qtx.DeleteAuthData(ctx, userUUID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Given a username and a password, validate the login attempt.
// Returns true if the username is valid, and the password matches the expected hash.
// Returns false otherwise.
//
// Fails and returns a non-nil error if:
// - The username does not exist in the database
func (database *DatabaseManager) ValidateLoginAttempt(ctx context.Context, username string, passwordAttempt string) (bool, error) {
	userDatum, err := database.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrOnFetchUserDoesNotExist
		}
		return false, err
	}
	authDatum, err := database.queries.GetAuthData(ctx, userDatum.Uuid)
	if err != nil {
		return false, err
	}

	if len(authDatum.Salt) != int(saltLen) {
		return false, errors.New("length of authDatum salt does not equal expected saltLen")
	}

	if len(authDatum.HashedPassword) != int(keyLen) {
		return false, errors.New("length of authDatum hashedPassword does not equal expected keyLen")
	}

	attemptHash := calculateHash(passwordAttempt, authDatum.Salt)
	return authDatum.HashedPassword == attemptHash, nil
}
