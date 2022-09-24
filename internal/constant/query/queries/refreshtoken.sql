-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (expires_at,
                            user_id,
                            scope,
                            redirect_uri,
                            client_id,
                            refresh_token,
                            code)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: RemoveRefreshTokenByCode :exec
DELETE
FROM refresh_tokens
WHERE code = $1;

-- name: RemoveRefreshToken :exec
DELETE
FROM refresh_tokens
WHERE refresh_token = $1;

-- name: GetRefreshTokenByUserIDAndClientID :one
SELECT *
FROM refresh_tokens
WHERE user_id = $1
  AND client_id = $2;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE refresh_token = $1;

-- name: GetAuthorizedClientsForUser :many
SELECT refresh_tokens.scope,
       refresh_tokens.expires_at,
       refresh_tokens.created_at,
       refresh_tokens.updated_at,
       clients.id,
       clients.name,
       clients.client_type,
       clients.logo_url
FROM refresh_tokens
         JOIN clients ON refresh_tokens.client_id = clients.id
WHERE user_id = $1;