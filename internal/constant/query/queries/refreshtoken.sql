-- name: SaveRefreshToken :one
INSERT INTO refreshtokens (
    expires_at,
    user_id,
    scope,
    redirect_uri,
    client_id,
    refreshtoken,
    code
) VALUES (
    $1, $2, $3, $4, $5,$6,$7
)
RETURNING *;

-- name: RemoveRefreshToken :exec
DELETE FROM refreshtokens WHERE code = $1;
