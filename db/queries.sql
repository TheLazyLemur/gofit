-- SQLite3 queries.
-- All queries should be idempotent.

-- name: Ping :one
SELECT 1;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?) RETURNING id;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token) VALUES (?, ?) RETURNING id;

-- name: JoinSessionByUserId :one
SELECT * FROM sessions JOIN users ON sessions.user_id = users.id WHERE sessions.token = ? LIMIT 1;
