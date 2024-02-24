-------------------------------------------------------------------------------
-- CREATE QUERIES

-- name: CreateAuthenticationData :one
INSERT INTO authenticationData(uuid, hashed_password, salt)
VALUES(?, ?, ?)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (uuid, username)
VALUES(?, ?)
RETURNING *;

-------------------------------------------------------------------------------
-- RETRIEVAL QUERIES

-- name: GetAuthData :one
SELECT * FROM authenticationData
WHERE uuid = ? 
LIMIT 1;

-- name: GetUserByUUID :one
SELECT * FROM users
WHERE uuid = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-------------------------------------------------------------------------------
-- UPDATE QUERIES

-- name: UpdateAuthenticationData :exec
UPDATE authenticationData
SET hashed_password = ?, salt = ?
WHERE uuid = ?;

-------------------------------------------------------------------------------
-- DELETE QUERIES

-- name: DeleteAuthData :exec
DELETE FROM authenticationData
WHERE uuid = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE uuid = ?;
