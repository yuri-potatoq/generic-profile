package enrollment

import (
	"context"
	"database/sql"
)

type Service interface {
	NewEnrollment(ctx context.Context) (id int, err error)
	GetEnrollmentState(ctx context.Context, id int) (EnrollmentState, error)
}

type service struct {
	r Repository
}

func NewEnrollmentService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) NewEnrollment(ctx context.Context) (id int, _ error) {
	return id, s.r.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		id, err = s.r.NewEnrollment(ctx, tx)
		return err
	})
}

func (s *service) GetEnrollmentState(ctx context.Context, id int) (EnrollmentState, error) {
	var stt EnrollmentState
	return stt, s.r.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		stt, err = s.r.GetEnrollmentState(ctx, tx, id)
		return err
	})
}
