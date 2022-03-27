package session

import (
	"context"
	"database/sql"
	"fmt"

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
	ID        int64  `json:"ID"`
	Hall_id   int64  `json:"hall_id,omitempty"`
	Movie_id  int64  `json:"movie_id,omitempty"`
	Starts_at string `json:"Starts_at"`
	VIP       bool   `json:"VIP"`
	Name      string `json:"Movie name"`
}

func (r *Resource) GID() int64 {
	return r.ID
}

// Create new entity in storage
func (r *Repository) Create(i internal.Identifiable, ctx context.Context) (internal.Identifiable, error) {
	var id int64

	session, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create session object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	err := sq.
		Insert("sessions").
		Columns("hall_id", "movie_id", "starts_at").
		Values(session.Hall_id, session.Movie_id, session.Starts_at).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRowContext(ctx).
		Scan(&id)

	if err != nil {
		r.Log.Info("Failed to run Create session query.",
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
		Select("sessions.id", "halls.vip", "movies.name", "starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		Where(sq.Eq{
			"sessions.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRowContext(ctx).
		Scan(&res.ID, &res.VIP, &res.Name, &res.Starts_at)

	if err == sql.ErrNoRows {

		return nil, nil
	}

	if err != nil {
		r.Log.Info("Failed to run Retrieve session query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return &res, nil
}

// Delete entity in storage
func (r *Repository) Delete(id int64, ctx context.Context) error {

	_, err := sq.
		Delete("sessions").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		ExecContext(ctx)

	if err != nil {
		r.Log.Info("Failed to run Delete session query.",
			zap.Error(err),
		)

		return internal.ErrInternalFailure
	}

	return nil
}

// RetrieveAll entity from storage
func (r *Repository) RetrieveAll(ctx context.Context) ([]internal.Identifiable, error) {

	rows, err := sq.
		Select("sessions.id", "halls.vip", "movies.name", "starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryContext(ctx)

	if err == sql.ErrNoRows {

		fmt.Println("NoRowsNoRowsNoRowsNoRowsNoRowsNoRowsNoRows")

		return nil, nil
	}

	if err != nil {
		r.Log.Info("Failed to run RetrieveAll sessions query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(&res.ID, &res.VIP, &res.Name, &res.Starts_at)

		if err != nil {
			r.Log.Info("Failed to scan rows into session structures.",
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

func (r *Repository) TimeValid(i internal.Identifiable, ctx context.Context) (bool, error) {
	session, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create session object.",
			zap.Bool("ok", ok),
		)

		return false, internal.ErrValidationFailed
	}

	res, err := sq.Select("movies.duration, sessions.starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Where("(?, movies.duration) OVERLAPS (sessions.starts_at , movies.duration) AND sessions.hall_id = ? AND sessions.movie_id = ?", session.Starts_at, session.Hall_id, session.Movie_id).
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)
	if err != nil {
		r.Log.Info("Failed to run query.",
			zap.Error(err),
		)

		return false, internal.ErrInternalFailure
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.Log.Info("Failed to count rows query.",
			zap.Error(err),
		)

		return false, internal.ErrInternalFailure
	}

	if rows == 0 {
		return true, nil
	}

	return false, internal.ErrValidationFailed
}
