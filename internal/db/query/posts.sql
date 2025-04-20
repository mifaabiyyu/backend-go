-- name: CreatePost :one
INSERT INTO posts (title, content, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING *;