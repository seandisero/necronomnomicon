-- name: CreateUser :one
INSERT INTO users (username, hashed_password, created_at, updated_at)
VALUES (?, ?, DATETIME('now'), DATETIME('now'))
RETURNING id, username, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByName :one
SELECT id, username, hashed_password, created_at, updated_at FROM users WHERE username = ?;

