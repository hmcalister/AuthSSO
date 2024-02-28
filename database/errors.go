package database

import "errors"

var (
	ErrOnCreateUserExists      error = errors.New("user exists in database")
	ErrOnFetchUserDoesNotExist error = errors.New("user does not exist in database")
)
