// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: refreshtoken.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getAuthorizedClientsForUser = `-- name: GetAuthorizedClientsForUser :many
SELECT refresh_tokens.scope,
       refresh_tokens.expires_at,
       refresh_tokens.created_at,
       refresh_tokens.updated_at,
       clients.id,
       clients.name,
       clients.client_type,
       clients.logo_url
FROM refresh_tokens
         JOIN clients ON refresh_tokens.client_id = clients.id
WHERE user_id = $1
  AND refresh_tokens.scope NOT ILIKE 'openid'
`

type GetAuthorizedClientsForUserRow struct {
	Scope      sql.NullString `json:"scope"`
	ExpiresAt  time.Time      `json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	ClientType string         `json:"client_type"`
	LogoUrl    string         `json:"logo_url"`
}

func (q *Queries) GetAuthorizedClientsForUser(ctx context.Context, userID uuid.UUID) ([]GetAuthorizedClientsForUserRow, error) {
	rows, err := q.db.Query(ctx, getAuthorizedClientsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAuthorizedClientsForUserRow
	for rows.Next() {
		var i GetAuthorizedClientsForUserRow
		if err := rows.Scan(
			&i.Scope,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID,
			&i.Name,
			&i.ClientType,
			&i.LogoUrl,
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

const getOpenIDAuthorizedClientsForUser = `-- name: GetOpenIDAuthorizedClientsForUser :many
SELECT refresh_tokens.scope,
       refresh_tokens.expires_at,
       refresh_tokens.created_at,
       refresh_tokens.updated_at,
       clients.id,
       clients.name,
       clients.client_type,
       clients.logo_url
FROM refresh_tokens
         JOIN clients ON refresh_tokens.client_id = clients.id
WHERE user_id = $1
  AND refresh_tokens.scope ILIKE '%openid%'
`

type GetOpenIDAuthorizedClientsForUserRow struct {
	Scope      sql.NullString `json:"scope"`
	ExpiresAt  time.Time      `json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	ClientType string         `json:"client_type"`
	LogoUrl    string         `json:"logo_url"`
}

func (q *Queries) GetOpenIDAuthorizedClientsForUser(ctx context.Context, userID uuid.UUID) ([]GetOpenIDAuthorizedClientsForUserRow, error) {
	rows, err := q.db.Query(ctx, getOpenIDAuthorizedClientsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOpenIDAuthorizedClientsForUserRow
	for rows.Next() {
		var i GetOpenIDAuthorizedClientsForUserRow
		if err := rows.Scan(
			&i.Scope,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID,
			&i.Name,
			&i.ClientType,
			&i.LogoUrl,
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

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, refresh_token, code, user_id, scope, redirect_uri, expires_at, client_id, created_at, updated_at
FROM refresh_tokens
WHERE refresh_token = $1
`

func (q *Queries) GetRefreshToken(ctx context.Context, refreshToken string) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshToken, refreshToken)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.RefreshToken,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.RedirectUri,
		&i.ExpiresAt,
		&i.ClientID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRefreshTokenByUserIDAndClientID = `-- name: GetRefreshTokenByUserIDAndClientID :one
SELECT id, refresh_token, code, user_id, scope, redirect_uri, expires_at, client_id, created_at, updated_at
FROM refresh_tokens
WHERE user_id = $1
  AND client_id = $2
`

type GetRefreshTokenByUserIDAndClientIDParams struct {
	UserID   uuid.UUID `json:"user_id"`
	ClientID uuid.UUID `json:"client_id"`
}

func (q *Queries) GetRefreshTokenByUserIDAndClientID(ctx context.Context, arg GetRefreshTokenByUserIDAndClientIDParams) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshTokenByUserIDAndClientID, arg.UserID, arg.ClientID)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.RefreshToken,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.RedirectUri,
		&i.ExpiresAt,
		&i.ClientID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeRefreshToken = `-- name: RemoveRefreshToken :exec
DELETE
FROM refresh_tokens
WHERE refresh_token = $1
`

func (q *Queries) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	_, err := q.db.Exec(ctx, removeRefreshToken, refreshToken)
	return err
}

const removeRefreshTokenByCode = `-- name: RemoveRefreshTokenByCode :exec
DELETE
FROM refresh_tokens
WHERE code = $1
`

func (q *Queries) RemoveRefreshTokenByCode(ctx context.Context, code string) error {
	_, err := q.db.Exec(ctx, removeRefreshTokenByCode, code)
	return err
}

const saveRefreshToken = `-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (expires_at,
                            user_id,
                            scope,
                            redirect_uri,
                            client_id,
                            refresh_token,
                            code)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, refresh_token, code, user_id, scope, redirect_uri, expires_at, client_id, created_at, updated_at
`

type SaveRefreshTokenParams struct {
	ExpiresAt    time.Time      `json:"expires_at"`
	UserID       uuid.UUID      `json:"user_id"`
	Scope        sql.NullString `json:"scope"`
	RedirectUri  sql.NullString `json:"redirect_uri"`
	ClientID     uuid.UUID      `json:"client_id"`
	RefreshToken string         `json:"refresh_token"`
	Code         string         `json:"code"`
}

func (q *Queries) SaveRefreshToken(ctx context.Context, arg SaveRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, saveRefreshToken,
		arg.ExpiresAt,
		arg.UserID,
		arg.Scope,
		arg.RedirectUri,
		arg.ClientID,
		arg.RefreshToken,
		arg.Code,
	)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.RefreshToken,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.RedirectUri,
		&i.ExpiresAt,
		&i.ClientID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateOAuthRefreshToken = `-- name: UpdateOAuthRefreshToken :one
UPDATE refresh_tokens
SET refresh_token = $1, updated_at = now()
WHERE refresh_token = $2
RETURNING id, refresh_token, code, user_id, scope, redirect_uri, expires_at, client_id, created_at, updated_at
`

type UpdateOAuthRefreshTokenParams struct {
	RefreshToken   string `json:"refresh_token"`
	RefreshToken_2 string `json:"refresh_token_2"`
}

func (q *Queries) UpdateOAuthRefreshToken(ctx context.Context, arg UpdateOAuthRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, updateOAuthRefreshToken, arg.RefreshToken, arg.RefreshToken_2)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.RefreshToken,
		&i.Code,
		&i.UserID,
		&i.Scope,
		&i.RedirectUri,
		&i.ExpiresAt,
		&i.ClientID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
