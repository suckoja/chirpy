-- name: CreateUser :one
INSERT INTO users (email)
VALUES (sqlc.arg(email))
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;