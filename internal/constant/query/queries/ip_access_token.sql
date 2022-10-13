-- name: SaveIPAccessToken :one
INSERT INTO ip_access_tokens (user_id, sub_id, ip_id, token, refresh_token)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetIPAccessTokenBySubAndIP :one
SELECT *
FROM ip_access_tokens
WHERE sub_id = $1
  AND ip_id = $2;

-- name: UpdateIPAccessToken :one
UPDATE ip_access_tokens
SET token         = sqlc.arg('token'),
    refresh_token = coalesce(sqlc.narg('refresh_token'), refresh_token)
WHERE sub_id = sqlc.arg('sub_id')
  AND ip_id = sqlc.arg('ip_id')
RETURNING *;