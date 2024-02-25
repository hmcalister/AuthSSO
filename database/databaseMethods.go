package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hmcalister/AuthSSO/database/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

// go:embed schema.sql
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

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := sqlc.New(db)

	return &DatabaseManager{
		db:      db,
		queries: queries,
	}, nil
}

// Given a username and password, attempt to register a new user
//
// Fails (and returns a non-nil error) if:
// - The salt fails to be generated
// - The username already exists in the database
// - The transaction to store both the new user data and the new auth data fails
//
// This method ensures that the new user data and auth data is create atomically, so
// a user cannot exist without auth data, and auth data cannot exist without a user
func (database *DatabaseManager) RegisterNewUser(username string, password string) error {
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

	ctx := context.Background()

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
func (database *DatabaseManager) ValidateAuthenticationAttempt(username string, passwordAttempt string) (bool, error) {
	ctx := context.Background()

	userDatum, err := database.queries.GetUserByUsername(ctx, username)
	if err != nil {
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
