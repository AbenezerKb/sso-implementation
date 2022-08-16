-- name: CreateClient :one
INSERT INTO clients (
    name,
    client_type,
    redirect_uris,
    scopes,
    secret,
    logo_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;
