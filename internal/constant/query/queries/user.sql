-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: GetUserByPhoneOrUserNameOrEmail :one
SELECT * FROM users WHERE phone = $1 OR user_name = $2 OR email = $3;

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
