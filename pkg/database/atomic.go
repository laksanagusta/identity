package database

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Manager interface {
	Atomic(
		ctx context.Context,
		callback func(ctx context.Context, tx DBTx) error,
	) error
}

type SQLTxManager struct {
	db DB
}

func NewManager(db DB) *SQLTxManager {
	return &SQLTxManager{db: db}
}

func (m *SQLTxManager) Atomic(ctx context.Context, callback func(ctx context.Context, tx DBTx) error) (rErr error) {
	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if rErr != nil {
			rErr = multierr.Combine(rErr, errors.WithStack(tx.Rollback()))
		}
	}()

	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				rErr = e
			} else {
				rErr = errors.Errorf("%s", rec)
			}
		}
	}()

	if err = callback(ctx, tx); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
