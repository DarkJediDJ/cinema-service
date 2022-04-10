package user_privileges

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
)

// Repository is a struct to store storage and logger connection
type Repository struct {
	DB  *sql.DB
	Log *zap.Logger
}

// Resource is a struct to store data about entity
type Resource struct {
	ID           int64  `json:"ID"`
	User_id      int64  `json:"User_id,omitempty"`
	Privilege_id int64  `json:"Privilege_id,omitempty"`
	Email        string `json:"Email"`
	Privilege    string `json:"Privilege"`
}

func (r *Resource) GID() int64 {
	return r.ID
}

// Create new entity in storage
func (r *Repository) Create(i internal.Identifiable, ctx context.Context) (internal.Identifiable, error) {
	var id int64

	user_privilege, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create user_privilege object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	err := sq.
		Insert("user_privileges").
		Columns("user_id", "privilege_id").
		Values(user_privilege.User_id, user_privilege.Privilege_id).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRowContext(ctx).
		Scan(&id)

	if err != nil {
		r.Log.Info("Failed to run Create user_privilege query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return r.Retrieve(id, ctx)
}

// Retrieve entity from storage
func (r *Repository) Retrieve(id int64, ctx context.Context) (internal.Identifiable, error) {
	var res Resource

	err := sq.
		Select("users.email", "user_privileges.id", "privileges.name").
		From("user_privileges").
		Join("users ON user_privileges.user_id = users.id").
		Join("privileges ON user_privileges.privilege_id = privileges.id").
		Where(sq.Eq{
			"user_privileges.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRowContext(ctx).
		Scan(&res.Email, &res.ID, &res.Privilege)

	if err == sql.ErrNoRows {

		return nil, nil
	}

	if err != nil {
		r.Log.Info("Failed to run Retrieve user_privilege query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return &res, nil
}

// Delete entity in storage
func (r *Repository) Delete(id int64, ctx context.Context) error {

	_, err := sq.
		Delete("user_privileges").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		ExecContext(ctx)

	if err != nil {
		r.Log.Info("Failed to run Delete user_privilege query.",
			zap.Error(err),
		)

		return internal.ErrInternalFailure
	}

	return nil
}

// RetrieveAll entity from storage
func (r *Repository) RetrieveAll(ctx context.Context) ([]internal.Identifiable, error) {

	rows, err := sq.
		Select("users.email", "user_privileges.id", "privileges.name").
		From("user_privileges").
		Join("users ON user_privileges.user_id = users.id").
		Join("privileges ON user_privileges.privilege_id = privileges.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).QueryContext(ctx)

	if err != nil {
		r.Log.Info("Failed to run RetrieveAll user_privilege query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(&res.Email, &res.ID, &res.Privilege)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			r.Log.Info("Failed to scan rows into user_privilege structures.",
				zap.Error(err),
			)

			return nil, internal.ErrInternalFailure
		}

		data = append(data, res)
	}

	var dataSlice []*Resource = data
	var interfaceSlice []internal.Identifiable = make([]internal.Identifiable, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}

	return interfaceSlice, nil
}
