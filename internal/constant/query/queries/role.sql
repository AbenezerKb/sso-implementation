-- name: AddRole :one
INSERT INTO roles (name, status)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteRole :one
DELETE
FROM roles
where name = $1
RETURNING *;

-- name: GetAllRoles :many
SELECT *
FROM roles;