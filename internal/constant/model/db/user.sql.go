// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
RETURNING id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at
`

type CreateUserParams struct {
	FirstName      string         `json:"first_name"`
	MiddleName     string         `json:"middle_name"`
	LastName       string         `json:"last_name"`
	Email          sql.NullString `json:"email"`
	Phone          string         `json:"phone"`
	UserName       string         `json:"user_name"`
	Password       string         `json:"password"`
	Gender         string         `json:"gender"`
	ProfilePicture sql.NullString `json:"profile_picture"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.MiddleName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.UserName,
		arg.Password,
		arg.Gender,
		arg.ProfilePicture,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const createUserWithID = `-- name: CreateUserWithID :one
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
RETURNING id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at
`

type CreateUserWithIDParams struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	MiddleName     string         `json:"middle_name"`
	LastName       string         `json:"last_name"`
	Email          sql.NullString `json:"email"`
	UserName       string         `json:"user_name"`
	Phone          string         `json:"phone"`
	Password       string         `json:"password"`
	Gender         string         `json:"gender"`
	ProfilePicture sql.NullString `json:"profile_picture"`
}

func (q *Queries) CreateUserWithID(ctx context.Context, arg CreateUserWithIDParams) (User, error) {
	row := q.db.QueryRow(ctx, createUserWithID,
		arg.ID,
		arg.FirstName,
		arg.MiddleName,
		arg.LastName,
		arg.Email,
		arg.UserName,
		arg.Phone,
		arg.Password,
		arg.Gender,
		arg.ProfilePicture,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, deleteUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at FROM users WHERE phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPhone, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByPhoneOrEmail = `-- name: GetUserByPhoneOrEmail :one
SELECT id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at FROM users WHERE phone = $1 OR email = $1
`

func (q *Queries) GetUserByPhoneOrEmail(ctx context.Context, phone string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPhoneOrEmail, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getUserStatus = `-- name: GetUserStatus :one
SELECT status FROM users WHERE id = $1
`

func (q *Queries) GetUserStatus(ctx context.Context, id uuid.UUID) (sql.NullString, error) {
	row := q.db.QueryRow(ctx, getUserStatus, id)
	var status sql.NullString
	err := row.Scan(&status)
	return status, err
}

const updatePhone = `-- name: UpdatePhone :exec
UPDATE users
SET phone = $1 WHERE phone = $2
`

type UpdatePhoneParams struct {
	NewPhone string `json:"new_phone"`
	OldPhone string `json:"old_phone"`
}

func (q *Queries) UpdatePhone(ctx context.Context, arg UpdatePhoneParams) error {
	_, err := q.db.Exec(ctx, updatePhone, arg.NewPhone, arg.OldPhone)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
 first_name = coalesce($1, first_name),
 middle_name = coalesce($2, middle_name),
 last_name = coalesce($3, last_name),
 email = coalesce($4, email),
 phone = coalesce($5, phone),
 user_name = coalesce($6, user_name),
 password = coalesce($7, password),
 gender = coalesce($8, gender),
 status = coalesce($9, status),
 profile_picture = coalesce($10)
WHERE id = $11
RETURNING id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at
`

type UpdateUserParams struct {
	FirstName      sql.NullString `json:"first_name"`
	MiddleName     sql.NullString `json:"middle_name"`
	LastName       sql.NullString `json:"last_name"`
	Email          sql.NullString `json:"email"`
	Phone          sql.NullString `json:"phone"`
	UserName       sql.NullString `json:"user_name"`
	Password       sql.NullString `json:"password"`
	Gender         sql.NullString `json:"gender"`
	Status         sql.NullString `json:"status"`
	ProfilePicture sql.NullString `json:"profile_picture"`
	ID             uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.FirstName,
		arg.MiddleName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.UserName,
		arg.Password,
		arg.Gender,
		arg.Status,
		arg.ProfilePicture,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserByID = `-- name: UpdateUserByID :one
UPDATE users
SET
 first_name = $2,
 middle_name = $3,
 last_name = $4,
 status = $5,
 phone = $6,
 profile_picture = $7
WHERE id = $1
RETURNING id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at
`

type UpdateUserByIDParams struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	MiddleName     string         `json:"middle_name"`
	LastName       string         `json:"last_name"`
	Status         sql.NullString `json:"status"`
	Phone          string         `json:"phone"`
	ProfilePicture sql.NullString `json:"profile_picture"`
}

func (q *Queries) UpdateUserByID(ctx context.Context, arg UpdateUserByIDParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserByID,
		arg.ID,
		arg.FirstName,
		arg.MiddleName,
		arg.LastName,
		arg.Status,
		arg.Phone,
		arg.ProfilePicture,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UserName,
		&i.Gender,
		&i.ProfilePicture,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
