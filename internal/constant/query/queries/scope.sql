-- name: CreateScope :one
INSERT INTO scopes (
    name,
    description,
    resource_server_name
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetScope :one
SELECT * FROM scopes WHERE name = $1;

-- name: DeleteScope :one
DELETE FROM scopes WHERE name = $1 RETURNING *;

-- name: GetScopesByResourceServerName :many
SELECT *
FROM scopes
WHERE resource_server_name = $1;