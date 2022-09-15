-- name: SaveInternalRefreshToken :one
INSERT INTO internalrefreshtokens (
    expires_at,
    user_id,
    refreshtoken
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: RemoveInternalRefreshToken :exec
DELETE FROM internalrefreshtokens WHERE refreshtoken =$1;

-- name: GetInternalRefreshToken :one
SELECT * FROM internalrefreshtokens WHERE refreshtoken = $1;

-- name: GetInternalRefreshTokenByUserID :one
SELECT * FROM internalrefreshtokens WHERE user_id = $1;

-- name: UpdateRefreshToken :one
Update internalrefreshtokens set expires_at = $2, refreshtoken= $3 WHERE id= $1 RETURNING *;