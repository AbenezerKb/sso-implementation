// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: auth_histories.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createAuthHistory = `-- name: CreateAuthHistory :one
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
RETURNING id, code, user_id, scope, status, redirect_uri, client_id, created_at
`

type CreateAuthHistoryParams struct {
	Code        string         `json:"code"`
	UserID      uuid.UUID      `json:"user_id"`
	Scope       sql.NullString `json:"scope"`
	RedirectUri sql.NullString `json:"redirect_uri"`
	ClientID    uuid.UUID      `json:"client_id"`
	Status      string         `json:"status"`
}

func (q *Queries) CreateAuthHistory(ctx context.Context, arg CreateAuthHistoryParams) (AuthHistory, error) {
	row := q.db.QueryRow(ctx, createAuthHistory,
		arg.Code,
		arg.UserID,
		arg.Scope,
		arg.RedirectUri,
		arg.ClientID,
		arg.Status,
	)
	var i AuthHistory
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.Status,
		&i.RedirectUri,
		&i.ClientID,
		&i.CreatedAt,
	)
	return i, err
}

const getAuthHistory = `-- name: GetAuthHistory :one
SELECT id, code, user_id, scope, status, redirect_uri, client_id, created_at FROM auth_histories WHERE code = $1
`

func (q *Queries) GetAuthHistory(ctx context.Context, code string) (AuthHistory, error) {
	row := q.db.QueryRow(ctx, getAuthHistory, code)
	var i AuthHistory
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.Status,
		&i.RedirectUri,
		&i.ClientID,
		&i.CreatedAt,
	)
	return i, err
}

const getLastAuthHistory = `-- name: GetLastAuthHistory :one
SELECT id, code, user_id, scope, status, redirect_uri, client_id, created_at FROM auth_histories WHERE user_id = $1 AND client_id = $2 ORDER BY created_at DESC LIMIT 1
`

type GetLastAuthHistoryParams struct {
	UserID   uuid.UUID `json:"user_id"`
	ClientID uuid.UUID `json:"client_id"`
}

func (q *Queries) GetLastAuthHistory(ctx context.Context, arg GetLastAuthHistoryParams) (AuthHistory, error) {
	row := q.db.QueryRow(ctx, getLastAuthHistory, arg.UserID, arg.ClientID)
	var i AuthHistory
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.Status,
		&i.RedirectUri,
		&i.ClientID,
		&i.CreatedAt,
	)
	return i, err
}
