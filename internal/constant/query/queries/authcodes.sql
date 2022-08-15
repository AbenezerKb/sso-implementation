-- name: CreateAuthCode :one
INSERT INTO authcodes (
    code,
    user_id,
    status,
    scope,
    redirect_uri,
    client_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAuthCode :one
SELECT * FROM authcodes WHERE code = $1;

-- name: DeleteAuthCode :one
DELETE FROM authcodes WHERE code = $1 RETURNING *;