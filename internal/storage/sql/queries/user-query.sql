-- name: GetUser :one
SELECT id, email, name, password, created_at, updated_at 
FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, name, password, created_at, updated_at 
FROM users 
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, name, password) 
VALUES ($1, $2, $3)
RETURNING id, email, name, password, created_at, updated_at;

-- name: UpdateUser :exec
UPDATE users 
SET name = $2, email = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at 
FROM users 
ORDER BY id;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);