-- name: CreateUser :one
INSERT INTO users (email)
VALUES (sqlc.arg(email))
RETURNING *;