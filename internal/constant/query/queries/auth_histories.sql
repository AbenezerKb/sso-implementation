-- name: CreateAuthHistory :one
INSERT INTO auth_histories (
    code,
    user_id,
    scope,
    redirect_uri,
    client_id,
    status
) VALUES (
    $1, $2, $3, $4, $5,$6
)
RETURNING *;

-- name: GetAuthHistory :one
SELECT * FROM auth_histories WHERE code = $1;

-- name: GetLastAuthHistory :one
SELECT * FROM auth_histories WHERE user_id = $1 AND client_id = $2 ORDER BY created_at DESC LIMIT 1;