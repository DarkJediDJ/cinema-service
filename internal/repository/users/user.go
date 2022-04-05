package user

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
)

// Repository is a struct to store DB and logger connection
type Repository struct {
	DB  *sql.DB
	Log *zap.Logger
}

// Resource is a struct to store data about entity
type Resource struct {
	ID       int64  `json:"ID"`
	EMail    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (r *Resource) GID() int64 {
	return r.ID
}

// Create new entity in storage
func (r *Repository) Create(i internal.Identifiable, ctx context.Context) error {
	user, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create use object.",
			zap.Bool("ok", ok),
		)

		return internal.ErrInternalFailure
	}

	_, err := sq.
		Insert("users").
		Columns("email", "password").
		Values(user.EMail, user.Password).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		ExecContext(ctx)

	if err != nil {
		r.Log.Info("Failed to run Create user query.",
			zap.Error(err),
		)

		return internal.ErrInternalFailure
	}

	return nil
}

// Retrieve entity from storage
func (r *Repository) Retrieve(email string, ctx context.Context) (internal.Identifiable, error) {
	var res Resource

	err := sq.
		Select("id", "email", "password").
		From("users").
		Where(sq.Eq{
			"email": email,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRowContext(ctx).
		Scan(&res.ID, &res.EMail, &res.Password)

	if err == sql.ErrNoRows {

		return nil, nil
	}

	if err != nil {
		r.Log.Info("Failed to run Retrieve user query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return &res, nil
}

// RetrievePrivileges entity from storage
func (r *Repository) RetrievePrivileges(id int64) ([]string, error) {

	var data []string

	rows, err := sq.
		Select("privileges.name").
		From("privileges").
		Join("user_privileges on user_privileges.privilege_id = privileges.id").
		Where(sq.Eq{
			"user_privileges.user_id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		Query()

	if err != nil {
		r.Log.Info("Failed to run Retrieve privileges query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}
	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			r.Log.Info("Failed to scan rows into privileges",
				zap.Error(err),
			)

			return nil, internal.ErrInternalFailure
		}

		data = append(data, name)
	}

	return data, nil
}

// RetrieveTickets entity from storage
func (r *Repository) RetrieveTickets(ticket int64, user int64) (bool, error) {
	var res string

	err := sq.
		Select("tickets.id").
		From("tickets").
		Where(sq.Eq{
			"tickets.id":      ticket,
			"tickets.user_id": user,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRow().
		Scan(&res)

	if err == sql.ErrNoRows {

		return false, nil
	}

	if err != nil {
		r.Log.Info("Failed to run Retrieve users tickets query.",
			zap.Error(err),
		)

		return false, internal.ErrInternalFailure
	}

	return true, nil
}
