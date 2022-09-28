-- name: CreateResourceServer :one
INSERT INTO resource_servers (name)
VALUES ($1)
RETURNING *;