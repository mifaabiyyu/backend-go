-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;