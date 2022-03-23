package session

import (
	"database/sql"

	"github.com/darkjedidj/cinema-service/internal"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
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
func (r *Repository) Create(i internal.Identifiable) (internal.Identifiable, error) {
	var id int

	session, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create session object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}
	query := sq.
		Insert("sessions").
		Columns("hall_id", "movie_id", "starts_at").
		Values(session.Hall_id, session.Movie_id, session.Starts_at).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)
	err := query.
		QueryRow().
		Scan(&id)

	if err != nil {
		r.Log.Info("Failed to run Create session query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return r.Retrieve(int64(id))
}

// Retrieve entity from storage
func (r *Repository) Retrieve(id int64) (internal.Identifiable, error) {
	var res Resource

	query := sq.
		Select("sessions.id", "halls.vip", "movies.name", "starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		Where(sq.Eq{
			"sessions.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err := query.
		QueryRow().
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
func (r *Repository) Delete(id int64) error {
	query := sq.
		Delete("sessions").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	_, err := query.
		Exec()

	if err != nil {
		r.Log.Info("Failed to run Delete session query.",
			zap.Error(err),
		)

		return internal.ErrInternalFailure
	}

	return nil
}

// RetrieveAll entity from storage
func (r *Repository) RetrieveAll() ([]internal.Identifiable, error) {
	query := sq.
		Select("sessions.id", "halls.vip", "movies.name", "starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
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
		if err == sql.ErrNoRows {
			return nil, nil
		}

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

func (r *Repository) TimeValid(i internal.Identifiable) (bool, error) {
	session, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create session object.",
			zap.Bool("ok", ok),
		)

		return false, internal.ErrValidationFailed
	}

	query := sq.Select("movies.duration, sessions.starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Where("(?, movies.duration) OVERLAPS (sessions.starts_at , movies.duration) AND sessions.hall_id = ? AND sessions.movie_id = ?", session.Starts_at, session.Hall_id, session.Movie_id).
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar)

	res, err := query.
		Exec()
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

	if err != nil {
		r.Log.Info("Failed to run time valid session query.",
			zap.Error(err),
		)

		return false, internal.ErrInternalFailure
	}

	return false, internal.ErrValidationFailed
}
