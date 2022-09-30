package persistencedb

import (
	"context"
	"sso/internal/constant/model/db"

	"github.com/jackc/pgx/v4"
)

func (p *PersistenceDB) SwapPhones(ctx context.Context, newPhone, oldPhone string) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := p.Queries.WithTx(tx)

	dummyPhone := newPhone + "d"
	err = qtx.UpdatePhone(ctx, db.UpdatePhoneParams{
		OldPhone: newPhone,
		NewPhone: dummyPhone,
	})
	if err != nil {
		return err
	}

	err = qtx.UpdatePhone(ctx, db.UpdatePhoneParams{
		OldPhone: oldPhone,
		NewPhone: newPhone,
	})

	if err != nil {
		return err
	}

	err = qtx.UpdatePhone(ctx, db.UpdatePhoneParams{
		OldPhone: dummyPhone,
		NewPhone: oldPhone,
	})

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
