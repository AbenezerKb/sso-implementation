-- name: SaveInternalRefreshToken :one
INSERT INTO internalrefreshtokens (
    expires_at,
    user_id,
    refresh_token,
    ip_address,
    user_agent
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: RemoveInternalRefreshToken :exec
DELETE FROM internalrefreshtokens WHERE refresh_token =$1;

-- name: GetInternalRefreshToken :one
SELECT * FROM internalrefreshtokens WHERE refresh_token = $1;

-- name: GetInternalRefreshTokensByUserID :many
SELECT * FROM internalrefreshtokens WHERE user_id = $1;

-- name: UpdateRefreshToken :one
Update internalrefreshtokens set expires_at = $2, refresh_token= $3 WHERE id= $1 RETURNING *;

-- name: RemoveInternalRefreshTokenByUserID :exec
DELETE FROM internalrefreshtokens WHERE id = $1;

-- name: UpdateInternalRefreshToken :one
UPDATE internalrefreshtokens SET refresh_token=$2, updated_at=now() WHERE refresh_token=$1 RETURNING *;