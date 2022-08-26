-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (
    expires_at,
    user_id,
    scope,
    redirect_uri,
    client_id,
    refresh_token,
    code
) VALUES (
    $1, $2, $3, $4, $5,$6,$7
)
RETURNING *;

-- name: RemoveRefreshTokenByCode :exec
DELETE FROM refresh_tokens WHERE code = $1;

-- name: RemoveRefreshToken :exec
DELETE FROM refresh_tokens WHERE refresh_token = $1;

-- name: CheckIfUserGrantedClient :one
SELECT * FROM refresh_tokens WHERE user_id = $1 AND client_id = $2;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE refresh_token = $1;

