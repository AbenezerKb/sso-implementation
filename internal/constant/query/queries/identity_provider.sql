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

-- name: UpdateIdentityProvider :one
UPDATE identity_providers
SET
    name = $2,
    logo_url = $3,
    client_id = $4,
    client_secret = $5,
    redirect_uri = $6,
    authorization_uri = $7,
    token_endpoint_url = $8,
    user_info_endpoint_url = $9
WHERE id = $1
RETURNING *;
