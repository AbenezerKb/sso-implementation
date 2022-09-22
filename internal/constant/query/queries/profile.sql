-- name: UpdateProfile :one
UPDATE users
SET
 first_name = $2,
 middle_name = $3,
 last_name = $4,
 gender = $5
WHERE id = $1
RETURNING *;