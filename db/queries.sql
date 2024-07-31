-- SQLite3 queries.
-- All queries should be idempotent.

-- name: Ping :one
SELECT 1;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?) RETURNING id;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token) VALUES (?, ?) RETURNING id;

-- name: GetUserByEmailAndPassword :one
SELECT * FROM users WHERE email = ? AND password_hash = ? LIMIT 1;

-- name: JoinSessionByUserId :one
SELECT * FROM sessions JOIN users ON sessions.user_id = users.id WHERE sessions.token = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?;

-- name: CreateUserWeight :exec
INSERT INTO user_weight (user_id, weight, created_at) VALUES (?, ?, ?);
