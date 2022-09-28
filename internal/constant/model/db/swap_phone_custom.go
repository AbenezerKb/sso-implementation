package db

import (
	"context"

	"github.com/jackc/pgx/v4"
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

func (p *PhoneSwap) SwapPhones(ctx context.Context, newPhone, oldPhone string) error {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := p.Queries.WithTx(tx)

	dummyPhone := newPhone + "d"
	err = qtx.UpdatePhone(ctx, UpdatePhoneParams{
		OldPhone: newPhone,
		NewPhone: dummyPhone,
	})
	if err != nil {
		return err
	}

	err = qtx.UpdatePhone(ctx, UpdatePhoneParams{
		OldPhone: oldPhone,
		NewPhone: newPhone,
	})

	if err != nil {
		return err
	}

	err = qtx.UpdatePhone(ctx, UpdatePhoneParams{
		OldPhone: dummyPhone,
		NewPhone: oldPhone,
	})

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
