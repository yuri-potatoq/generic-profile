package enrollment

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/yuri-potatoq/generic-profile/infra/db"
)

type Repository interface {
	db.TxManager
	NewEnrollment(ctx context.Context, tx *sql.Tx) (int, error)
	GetEnrollmentState(ctx context.Context, tx *sql.Tx, id int) (EnrollmentState, error)
}

type repository struct {
	db.TxManager
}

func NewEnrollmentRepository(d *sqlx.DB) Repository {
	return &repository{
		TxManager: db.NewTxManager(d),
	}
}

func (r *repository) NewEnrollment(ctx context.Context, tx *sql.Tx) (int, error) {
	rs, err := tx.ExecContext(ctx, "INSERT INTO enrollments DEFAULT VALUES RETURNING ID;")
	id, err := rs.LastInsertId()
	return int(id), err
}

func (r *repository) CheckEnrollment(ctx context.Context, tx *sql.Tx, id int) (bool, error) {
	var total int
	if err := tx.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM enrollments WHERE ID = $1;", id,
	).Scan(&total); err != nil {
		return false, err
	}
	return total > 0, nil
}

func (r *repository) GetEnrollmentState(ctx context.Context, tx *sql.Tx, id int) (EnrollmentState, error) {

	return EnrollmentState{}, nil
}
