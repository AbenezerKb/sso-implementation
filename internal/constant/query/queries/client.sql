-- name: CreateClient :one
INSERT INTO clients (
    name,
    client_type,
    redirect_uris,
    scopes,
    secret,
    logo_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: DeleteClient :one
DELETE FROM clients WHERE id = $1 RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients WHERE id = $1;

-- name: UpdateClient :one
UPDATE clients
SET
 name = coalesce(sqlc.narg('name'), name),
 client_type = coalesce(sqlc.narg('client_type'), client_type),
 redirect_uris = coalesce(sqlc.narg('redirect_uris'), redirect_uris),
 scopes = coalesce(sqlc.narg('scopes'), scopes),
 secret = coalesce(sqlc.narg('secret'), secret),
 logo_url = coalesce(sqlc.narg('logo_url'), logo_url),
 status = coalesce(sqlc.narg('status'), status)
WHERE id = sqlc.arg('id')
RETURNING *;
