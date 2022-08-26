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

-- name: RemoveRefreshTokenByCode :exec
DELETE FROM refreshtokens WHERE code = $1;

-- name: RemoveRefreshToken :exec
DELETE FROM refreshtokens WHERE refreshtoken = $1;

-- name: CheckIfUserGrantedClient :one
SELECT * FROM refreshtokens WHERE user_id = $1 AND client_id = $2;

-- name: GetRefreshToken :one
SELECT * FROM refreshtokens WHERE refreshtoken = $1;

