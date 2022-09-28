-- name: CreateResourceServer :one
INSERT INTO resource_servers (name)
VALUES ($1)
RETURNING *;

-- name: DeleteResourceServer :one
DELETE
FROM resource_servers
WHERE id = $1
RETURNING *;