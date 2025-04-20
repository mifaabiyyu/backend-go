-- name: CreateUser :one
INSERT INTO users (email, username, full_name, password, role_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
JOIN roles ON users.role_id = roles.id
LIMIT $1
OFFSET $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
  set email = $2, 
  username = $3, 
  full_name = $4,
  updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetUserWithRole :one
SELECT 
    users.id,
    users.email,
    users.username,
    users.full_name,
    users.password,
    users.verified,
    users.verified_at,
    users.created_at,
    users.updated_at,
    users.role_id,
    roles.* AS role
FROM users
JOIN roles ON users.role_id = roles.id
WHERE users.id = $1;

-- name: GetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;