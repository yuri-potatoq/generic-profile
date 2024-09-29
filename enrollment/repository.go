package enrollment

import (
	"context"
	"database/sql"
)

type Repository interface {
	NewEnrollment(ctx context.Context, )
}

type repository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) Repository {
	return repository{
		db: db,
	}
}
