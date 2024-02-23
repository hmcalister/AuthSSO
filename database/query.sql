-- name: CreateAuthenticationData :one
INSERT INTO authenticationData(uuid, hashedPassword, salt)
VALUES(?, ?, ?)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (uuid, username)
VALUES(?, ?)
RETURNING *;

-- name: GetAuthData :one
SELECT * FROM authenticationData
WHERE uuid = ? 
LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE uuid = ? LIMIT 1;

-- name: UpdateAuthenticationData :exec
UPDATE authenticationData
SET hashedPassword = ?, salt = ?
WHERE uuid = ?;

-- name: DeleteAuthData :exec
DELETE FROM authenticationData
WHERE uuid = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE uuid = ?;
