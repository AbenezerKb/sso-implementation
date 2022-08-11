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

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, first_name, middle_name, last_name, email, phone, password, user_name, gender, profile_picture, status, created_at FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
