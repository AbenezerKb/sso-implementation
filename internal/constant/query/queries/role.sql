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

-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE name = $1;

-- name: GetRoleStatus :one
SELECT status
FROM roles
WHERE name = $1;