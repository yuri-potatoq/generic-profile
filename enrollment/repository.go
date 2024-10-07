package enrollment

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yuri-potatoq/generic-profile/infra/db"
)

type Repository interface {
	db.TxManager
	NewEnrollment(ctx context.Context) db.Querier[int]
	CheckEnrollment(ctx context.Context, id int) db.Querier[bool]
	GetModalities(ctx context.Context, id int) db.Querier[[]Modalities]
	GetChildProfile(ctx context.Context, id int) db.Querier[[]ChildProfile]
	GetChildParent(ctx context.Context, id int) db.Querier[[]ChildParent]
	GetAddress(ctx context.Context, id int) db.Querier[[]Address]
	GetShifts(ctx context.Context, id int) db.Querier[[]EnrollmentShift]
	GetTerms(ctx context.Context, id int) db.Querier[bool]
	UpdateTerm(ctx context.Context, id int, term bool) db.Executer
	UpdateChildProfile(ctx context.Context, id int, p ChildProfile) db.Executer
	UpdateModalities(ctx context.Context, id int, modalities []Modalities) db.Executer
	UpdateShift(ctx context.Context, id int, shift EnrollmentShift) db.Executer
	UpdateAddress(ctx context.Context, id int, addr Address) db.Executer
	UpdateChildParent(ctx context.Context, id int, parent ChildParent) db.Executer
}

type repository struct {
	db.TxManager
}

func NewEnrollmentRepository(d *sqlx.DB) Repository {
	return &repository{
		TxManager: db.NewTxManager(d),
	}
}

func (r *repository) NewEnrollment(ctx context.Context) db.Querier[int] {
	return db.NewQuerier[int](ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (int, error) {
		rs, err := tx.ExecContext(iCtx, "INSERT INTO enrollments DEFAULT VALUES RETURNING ID;")
		id, err := rs.LastInsertId()
		return int(id), err
	})
}

func (r *repository) CheckEnrollment(ctx context.Context, id int) db.Querier[bool] {
	return db.NewQuerier[bool](ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (bool, error) {
		var total int
		if err := tx.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM enrollments WHERE ID = $1;",
			id,
		).Scan(&total); err != nil {
			return false, err
		}
		return total > 0, nil
	})
}

func (r *repository) GetAddress(ctx context.Context, id int) db.Querier[[]Address] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (addrs []Address, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT
				a.city
				,a.house_number
				,a.state
				,a.street
				,a.zipcode
			FROM addresses a
			INNER JOIN enrollments_addresses ea
				ON ea.addresses_fk = a.ID
			WHERE ea.enrollment_fk =  $1
			ORDER BY a.ID DESC`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return addrs, sqlx.StructScan(rows, &addrs)
	})
}

func (r *repository) GetChildProfile(ctx context.Context, id int) db.Querier[[]ChildProfile] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (profiles []ChildProfile, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT 
				full_name
				,birthdate
				,gender
				,medical_info
			FROM child_profile cp
			JOIN enrollments_child ec
				ON ec.child_profile_fk = cp.ID
			WHERE ec.enrollment_fk = $1
			ORDER BY cp.ID DESC`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return profiles, sqlx.StructScan(rows, &profiles)
	})
}

func (r *repository) GetChildParent(ctx context.Context, id int) db.Querier[[]ChildParent] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (parents []ChildParent, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT
				p.email
				,p.full_name
				,p.phone_number
			FROM parents p
			JOIN child_parents ec 
				ON ec.parents_fk = p.ID
			WHERE ec.enrollment_fk = $1
			ORDER BY p.ID DESC;`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return parents, sqlx.StructScan(rows, &parents)
	})
}

func (r *repository) GetModalities(ctx context.Context, id int) db.Querier[[]Modalities] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (modalities []Modalities, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT m.name FROM modalities m
			LEFT JOIN student_modalities sm 
			    ON m.ID = sm.modalities_fk
			WHERE sm.enrollment_fk = $1`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var modality Modalities
			err = rows.Scan(&modality)
			if err != nil {
				return nil, err
			}
			modalities = append(modalities, modality)
		}
		return modalities, nil
	})
}

func (r *repository) GetShifts(ctx context.Context, id int) db.Querier[[]EnrollmentShift] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (shifts []EnrollmentShift, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT es.name FROM enrollments_shift es
			LEFT JOIN enrollments_shifts ess
			    ON es.ID = ess.enrollments_shift_fk
			WHERE ess.enrollment_fk = $1;`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var shift EnrollmentShift
			err = rows.Scan(&shift)
			if err != nil {
				return nil, err
			}
			shifts = append(shifts, shift)
		}
		return shifts, nil
	})
}

func (r *repository) GetTerms(ctx context.Context, id int) db.Querier[bool] {
	return db.NewQuerier(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) (term bool, err error) {
		rows, err := tx.QueryContext(iCtx, `
			SELECT count(*) FROM enrollments_terms 
				WHERE terms_agreement = true AND enrollment_fk = $1;`, id)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		for rows.Next() {
			var term bool
			return term, rows.Scan(&term)
		}
		return false, nil
	})
}

func (r *repository) UpdateTerm(ctx context.Context, id int, term bool) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `
			INSERT INTO enrollments_terms(enrollment_fk, terms_agreement) 
				VALUES ($1, $2);`,
			id, term,
		)
		return err
	})
}

func (r *repository) UpdateChildProfile(ctx context.Context, id int, p ChildProfile) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `
			CALL insert_child_profile($1, $2, $3, $4, $5);
		`,
			id, p.FullName, p.Birthdate, p.Gender, p.MedicalInfo,
		)
		return err
	})
}

func (r *repository) UpdateShift(ctx context.Context, id int, shift EnrollmentShift) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `
				CALL insert_shift($1, $2);
			`,
			id, shift,
		)
		return err
	})
}
func (r *repository) UpdateModalities(ctx context.Context, id int, modalities []Modalities) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `
				CALL insert_modality($1, $2::TEXT[]);
			`,
			id, pq.Array(modalities),
		)
		return err
	})
}

func (r *repository) UpdateAddress(ctx context.Context, id int, addr Address) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `				
				CALL insert_address($1, $2, $3, $4, $5);
			`,
			id, addr.ZipCode, addr.State, addr.City, addr.HouseNumber,
		)
		return err
	})
}

func (r *repository) UpdateChildParent(ctx context.Context, id int, parent ChildParent) db.Executer {
	return db.NewExecuter(ctx, r.TxManager, func(iCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(iCtx, `				
				CALL insert_child_parent($1, $2, $3, $4);
			`,
			id, parent.FullName, parent.PhoneNumber, parent.Email,
		)
		return err
	})
}
