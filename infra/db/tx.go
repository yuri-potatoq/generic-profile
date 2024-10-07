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

type Querier[T any] interface {
	Atomic() (T, error)
	Tx(tx *sql.Tx) (T, error)
}

type querier[T any] struct {
	TxManager
	ctx    context.Context
	queryF func(ctx context.Context, tx *sql.Tx) (T, error)
}

func (q *querier[T]) Atomic() (T, error) {
	var rs T
	return rs, q.TxManager.WithTx(q.ctx, func(tx *sql.Tx) error {
		var err error
		rs, err = q.queryF(q.ctx, tx)
		return err
	})
}

func (q *querier[T]) Tx(tx *sql.Tx) (T, error) {
	return q.queryF(q.ctx, tx)
}

func NewQuerier[T any](
	ctx context.Context,
	tx TxManager,
	f func(ctx context.Context, tx *sql.Tx,
) (T, error)) Querier[T] {
	return &querier[T]{
		ctx:       ctx,
		queryF:    f,
		TxManager: tx,
	}
}

type Executer interface {
	Atomic() error
	Tx(tx *sql.Tx) error
}

type executer struct {
	TxManager
	ctx   context.Context
	execF func(ctx context.Context, tx *sql.Tx) error
}

func NewExecuter(
	ctx context.Context,
	tx TxManager,
	f func(ctx context.Context, tx *sql.Tx) error,
) Executer {
	return &executer{
		TxManager: tx,
		ctx:       ctx,
		execF:     f,
	}
}

func (q *executer) Atomic() error {
	return q.TxManager.WithTx(q.ctx, func(tx *sql.Tx) error {
		return q.execF(q.ctx, tx)
	})
}

func (q *executer) Tx(tx *sql.Tx) error {
	return q.execF(q.ctx, tx)
}
