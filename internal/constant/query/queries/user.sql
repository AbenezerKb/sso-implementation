-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: GetUserStatus :one
SELECT status FROM users WHERE id = $1;

-- name: GetUserByPhoneOrEmail :one
SELECT * FROM users WHERE phone = $1 OR email = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
    first_name,
    middle_name,
    last_name,
    email,
    phone,
    user_name,
    password,
    gender,
    profile_picture
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
 first_name = coalesce(sqlc.narg('first_name'), first_name),
 middle_name = coalesce(sqlc.narg('middle_name'), middle_name),
 last_name = coalesce(sqlc.narg('last_name'), last_name),
 email = coalesce(sqlc.narg('email'), email),
 phone = coalesce(sqlc.narg('phone'), phone),
 user_name = coalesce(sqlc.narg('user_name'), user_name),
 password = coalesce(sqlc.narg('password'), password),
 gender = coalesce(sqlc.narg('gender'), gender),
 status = coalesce(sqlc.narg('status'), status),
 profile_picture = coalesce(sqlc.narg('profile_picture'))
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: UpdateUserByID :one
UPDATE users
SET
 first_name = $2,
 middle_name = $3,
 last_name = $4,
 status = $5,
 phone = $6,
 profile_picture = $7
WHERE id = $1
RETURNING *;

-- name: UpdatePhone :exec
UPDATE users
SET phone = sqlc.arg('new_phone') WHERE phone = sqlc.arg('old_phone');

-- name: CreateUserWithID :one
INSERT INTO users (
    id,
    first_name,
    middle_name,
    last_name,
    email,
    user_name,
    phone,
    password,
    gender,
    profile_picture
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;