-- name: GetAllUsers :many
SELECT * FROM users;

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