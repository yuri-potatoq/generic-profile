package enrollment

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/yuri-potatoq/generic-profile/infra/db"
)

type Service interface {
	GetEnrollmentState(ctx context.Context, id int) (EnrollmentState, error)
	BulkUpdate(ctx context.Context, partialUpdate PartialUpdate) (EnrollmentState, error)
}

type service struct {
	r Repository
}

func NewEnrollmentService(r Repository) Service {
	return &service{
		r: r,
	}
}

func firstOrDefault[T any](xs []T) T {
	if len(xs) > 1 {
		return xs[0]
	}
	return *new(T)
}

func (s *service) GetEnrollmentState(ctx context.Context, enrollmentId int) (EnrollmentState, error) {
	var stt EnrollmentState
	return stt, s.r.WithTx(ctx, func(tx *sql.Tx) (err error) {
		if exists, err := s.r.CheckEnrollment(ctx, enrollmentId).Tx(tx); err != nil {
			return err
		} else if !exists {
			return fmt.Errorf("can't find enrollment with id (%d)", enrollmentId)
		}

		profiles, err := s.r.GetChildProfile(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}
		stt.ChildProfile = firstOrDefault(profiles)

		parent, err := s.r.GetChildParent(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}
		stt.ChildParent = firstOrDefault(parent)

		addr, err := s.r.GetAddress(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}
		stt.Address = firstOrDefault(addr)

		shifts, err := s.r.GetShifts(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}
		stt.EnrollmentShift = firstOrDefault(shifts)

		stt.Terms, err = s.r.GetTerms(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}

		stt.Modalities, err = s.r.GetModalities(ctx, enrollmentId).Tx(tx)
		return err
	})
}

type PartialUpdate struct {
	ops []partialOperation
	ID  *int
}

type (
	partialOperation func(
		ctx context.Context,
		s *service,
		enrollmentId int,
	) (db.Executer, error)

	// Partial update option is a save layer for partialOperation closures.
	// Which means, we can avoid expose inner types of service layer.
	PartialUpdateOpt func() partialOperation
)

func NewPartialUpdate(opts ...PartialUpdateOpt) PartialUpdate {
	var partial PartialUpdate
	for _, opt := range opts {
		partial.ops = append(partial.ops, opt())
	}
	return partial
}

func UpdateWithChildProfile(v ChildProfile) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateChildProfile(ctx, enrollmentId, v), nil
		}
	}
}

func UpdateWithChildParent(v ChildParent) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateChildParent(ctx, enrollmentId, v), nil
		}
	}
}

func UpdateWithAddress(v Address) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateAddress(ctx, enrollmentId, v), nil
		}
	}
}

func UpdateWithModalities(v []Modalities) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateModalities(ctx, enrollmentId, v), nil
		}
	}
}

func UpdateWithShift(v EnrollmentShift) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateShift(ctx, enrollmentId, v), nil
		}
	}
}

func UpdateWithTerm(v bool) PartialUpdateOpt {
	return func() partialOperation {
		return func(ctx context.Context, s *service, enrollmentId int) (db.Executer, error) {
			return s.r.UpdateTerm(ctx, enrollmentId, v), nil
		}
	}
}

//func (par *PartialUpdate) IterOps(ctx context.Context, s *service, tx *sql.Tx) iter.Seq[error] {
//	return func(yield func(error) bool) {
//		for _, op := range par.ops {
//			exec := op(ctx, s)
//			err := exec.Tx(tx)
//			if !yield(err) {
//				return
//			}
//		}
//	}
//}

func (s *service) BulkUpdate(ctx context.Context, partialUpdate PartialUpdate) (EnrollmentState, error) {
	var stt EnrollmentState
	err := s.r.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		var enrollmentId int
		if partialUpdate.ID == nil {
			enrollmentId = -1
		} else {
			enrollmentId = *partialUpdate.ID
		}
		exists, err := s.r.CheckEnrollment(ctx, enrollmentId).Tx(tx)
		if err != nil {
			return err
		}
		if !exists {
			enrollmentId, err = s.r.NewEnrollment(ctx).Tx(tx)
			if err != nil {
				return err
			}
		}

		for _, op := range partialUpdate.ops {
			exec, err := op(ctx, s, enrollmentId)
			if err != nil {
				// TODO: invalid operation error
				return err
			}

			if err := exec.Tx(tx); err != nil {
				return err
			}
		}

		return err
	})

	return stt, err
}
