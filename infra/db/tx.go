package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"runtime/debug"
)

type TxManager interface {
	WithTx(ctx context.Context, op func(tx *sql.Tx) error) error
}

type commonManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) TxManager {
	return &commonManager{db}
}

func (m *commonManager) WithTx(ctx context.Context, op func(tx *sql.Tx) error) error {
	var err error
	c, err := m.db.Conn(ctx)
	if err != nil {
		return err
	}

	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			debug.PrintStack()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return op(tx)
}
