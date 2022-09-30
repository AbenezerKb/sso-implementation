package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type PhoneSwap struct {
	*Queries
	db *pgxpool.Pool
}

func NewPhoneSwap(db *pgxpool.Pool) PhoneSwap {
	return PhoneSwap{
		db:      db,
		Queries: New(db),
	}
}
