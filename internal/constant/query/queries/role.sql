-- name: AddRole :one
INSERT INTO roles (name)
VALUES ($1)
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

-- name: UpdateRoleStatus :one
UPDATE roles
SET status = $2
WHERE name = $1
RETURNING *;