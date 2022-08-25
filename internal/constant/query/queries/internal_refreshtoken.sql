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
DELETE FROM internalrefreshtokens WHERE user_id = $1;
