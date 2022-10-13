-- name: CreateIdentityProvider :one
INSERT INTO identity_providers (name, logo_url, client_id, client_secret, redirect_uri, authorization_uri,
                                token_endpoint_url, user_info_endpoint_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteIdentityProvider :one
DELETE
FROM identity_providers
WHERE id = $1
RETURNING *;

-- name: GetIdentityProvider :one
SELECT *
FROM identity_providers
WHERE id = $1;