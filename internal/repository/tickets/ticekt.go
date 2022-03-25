package tickets

import (
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
	Starts_at  string `json:"Starts_at"`
	Price      float64
	Seat       int
	ID         int64
	Title      string
	User_ID    int
	Session_ID int
}

func (r *Resource) GID() int64 {
	return r.ID
}

// Create new entity in storage
func (r *Repository) Create(i internal.Identifiable) (internal.Identifiable, error) {
	var id int

	ticket, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create ticket object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	err := sq.
		Insert("tickets").
		Columns("user_id", "price", "session_id").
		Values(ticket.User_ID, ticket.Price, ticket.Session_ID).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		r.Log.Info("Failed to run Create tickets query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	return r.Retrieve(int64(id))
}

// Retrieve entity from storage
func (r *Repository) Retrieve(id int64) (internal.Identifiable, error) {
	var res Resource

	err := sq.
		Select("tickets.id", "user_id", "price", "session_id", "movies.name", "halls.seats", "sessions.starts_at", "tickets.price").
		From("tickets").
		Join("sessions ON tickets.session_id = sessions.id").
		Join("movies ON sessions.movies_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		Where(sq.Eq{
			"tickets.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRow().
		Scan(&res.ID, &res.User_ID, &res.Price, &res.Session_ID, &res.Title, &res.Seat, &res.Starts_at, &res.Price)

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

	_, err := sq.
		Delete("sessions").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
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

	rows, err := sq.
		Select("sessions.id", "halls.vip", "movies.name", "starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).Query()
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

	res, err := sq.Select("movies.duration, sessions.starts_at").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Where("(?, movies.duration) OVERLAPS (sessions.starts_at , movies.duration) AND sessions.hall_id = ? AND sessions.movie_id = ?", session.Starts_at, session.Hall_id, session.Movie_id).
		RunWith(r.DB).
		PlaceholderFormat(sq.Dollar).
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
