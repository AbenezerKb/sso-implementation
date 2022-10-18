// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: client.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createClient = `-- name: CreateClient :one
INSERT INTO clients (
    name,
    client_type,
    redirect_uris,
    scopes,
    secret,
    logo_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, name, client_type, redirect_uris, scopes, secret, logo_url, status, created_at, first_party
`

type CreateClientParams struct {
	Name         string `json:"name"`
	ClientType   string `json:"client_type"`
	RedirectUris string `json:"redirect_uris"`
	Scopes       string `json:"scopes"`
	Secret       string `json:"secret"`
	LogoUrl      string `json:"logo_url"`
}

func (q *Queries) CreateClient(ctx context.Context, arg CreateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, createClient,
		arg.Name,
		arg.ClientType,
		arg.RedirectUris,
		arg.Scopes,
		arg.Secret,
		arg.LogoUrl,
	)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientType,
		&i.RedirectUris,
		&i.Scopes,
		&i.Secret,
		&i.LogoUrl,
		&i.Status,
		&i.CreatedAt,
		&i.FirstParty,
	)
	return i, err
}

const deleteClient = `-- name: DeleteClient :one
DELETE FROM clients WHERE id = $1 RETURNING id, name, client_type, redirect_uris, scopes, secret, logo_url, status, created_at, first_party
`

func (q *Queries) DeleteClient(ctx context.Context, id uuid.UUID) (Client, error) {
	row := q.db.QueryRow(ctx, deleteClient, id)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientType,
		&i.RedirectUris,
		&i.Scopes,
		&i.Secret,
		&i.LogoUrl,
		&i.Status,
		&i.CreatedAt,
		&i.FirstParty,
	)
	return i, err
}

const getClientByID = `-- name: GetClientByID :one
SELECT id, name, client_type, redirect_uris, scopes, secret, logo_url, status, created_at, first_party FROM clients WHERE id = $1
`

func (q *Queries) GetClientByID(ctx context.Context, id uuid.UUID) (Client, error) {
	row := q.db.QueryRow(ctx, getClientByID, id)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientType,
		&i.RedirectUris,
		&i.Scopes,
		&i.Secret,
		&i.LogoUrl,
		&i.Status,
		&i.CreatedAt,
		&i.FirstParty,
	)
	return i, err
}

const updateClient = `-- name: UpdateClient :one
UPDATE clients
SET
 name = coalesce($1, name),
 client_type = coalesce($2, client_type),
 redirect_uris = coalesce($3, redirect_uris),
 scopes = coalesce($4, scopes),
 secret = coalesce($5, secret),
 logo_url = coalesce($6, logo_url),
 status = coalesce($7, status)
WHERE id = $8
RETURNING id, name, client_type, redirect_uris, scopes, secret, logo_url, status, created_at, first_party
`

type UpdateClientParams struct {
	Name         sql.NullString `json:"name"`
	ClientType   sql.NullString `json:"client_type"`
	RedirectUris sql.NullString `json:"redirect_uris"`
	Scopes       sql.NullString `json:"scopes"`
	Secret       sql.NullString `json:"secret"`
	LogoUrl      sql.NullString `json:"logo_url"`
	Status       sql.NullString `json:"status"`
	ID           uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateClient(ctx context.Context, arg UpdateClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, updateClient,
		arg.Name,
		arg.ClientType,
		arg.RedirectUris,
		arg.Scopes,
		arg.Secret,
		arg.LogoUrl,
		arg.Status,
		arg.ID,
	)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientType,
		&i.RedirectUris,
		&i.Scopes,
		&i.Secret,
		&i.LogoUrl,
		&i.Status,
		&i.CreatedAt,
		&i.FirstParty,
	)
	return i, err
}

const updateEntireClient = `-- name: UpdateEntireClient :one
UPDATE clients
SET
 name = $2,
 client_type = $3,
 redirect_uris = $4,
 scopes = $5,
 logo_url = $6
WHERE id = $1
RETURNING id, name, client_type, redirect_uris, scopes, secret, logo_url, status, created_at, first_party
`

type UpdateEntireClientParams struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ClientType   string    `json:"client_type"`
	RedirectUris string    `json:"redirect_uris"`
	Scopes       string    `json:"scopes"`
	LogoUrl      string    `json:"logo_url"`
}

func (q *Queries) UpdateEntireClient(ctx context.Context, arg UpdateEntireClientParams) (Client, error) {
	row := q.db.QueryRow(ctx, updateEntireClient,
		arg.ID,
		arg.Name,
		arg.ClientType,
		arg.RedirectUris,
		arg.Scopes,
		arg.LogoUrl,
	)
	var i Client
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClientType,
		&i.RedirectUris,
		&i.Scopes,
		&i.Secret,
		&i.LogoUrl,
		&i.Status,
		&i.CreatedAt,
		&i.FirstParty,
	)
	return i, err
}
