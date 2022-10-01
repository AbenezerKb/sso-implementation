// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: resource_server.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createResourceServer = `-- name: CreateResourceServer :one
INSERT INTO resource_servers (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at
`

func (q *Queries) CreateResourceServer(ctx context.Context, name string) (ResourceServer, error) {
	row := q.db.QueryRow(ctx, createResourceServer, name)
	var i ResourceServer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteResourceServer = `-- name: DeleteResourceServer :one
DELETE
FROM resource_servers
WHERE id = $1
RETURNING id, name, created_at, updated_at
`

func (q *Queries) DeleteResourceServer(ctx context.Context, id uuid.UUID) (ResourceServer, error) {
	row := q.db.QueryRow(ctx, deleteResourceServer, id)
	var i ResourceServer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getResourceServerByName = `-- name: GetResourceServerByName :one
SELECT id, name, created_at, updated_at
FROM resource_servers
WHERE name = $1
`

func (q *Queries) GetResourceServerByName(ctx context.Context, name string) (ResourceServer, error) {
	row := q.db.QueryRow(ctx, getResourceServerByName, name)
	var i ResourceServer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
