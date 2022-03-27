package tickets

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
	Starts_at  string `json:"Starts_at"`
	Price      float64
	Seat       int64
	ID         int64
	Title      string
	User_ID    int64
	Session_ID int64
}

func (r *Resource) GID() int64 {
	return r.ID
}

// Create new entity in storage
func (r *Repository) Create(ctx context.Context, i internal.Identifiable, tx *sql.Tx) (int64, error) {
	var id int

	ticket, ok := i.(*Resource)
	if !ok {
		r.Log.Info("Failed to create ticket object.",
			zap.Bool("ok", ok),
		)

		return 0, internal.ErrInternalFailure
	}

	err := sq.
		Insert("tickets").
		Columns("user_id", "price", "session_id", "seat").
		Values(ticket.User_ID, ticket.Price, ticket.Session_ID, ticket.Seat).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		QueryRow().
		Scan(&id)

	if err != nil {
		r.Log.Info("Failed to run Create tickets query.",
			zap.Error(err),
		)

		return 0, internal.ErrInternalFailure
	}

	return int64(id), nil
}

// Retrieve entity from storage
func (r *Repository) Retrieve(id int64, ctx context.Context) (internal.Identifiable, error) {
	var res Resource

	err := sq.
		Select("tickets.id", "user_id", "price", "session_id", "movies.name", "tickets.seat", "sessions.starts_at").
		From("tickets").
		Join("sessions ON tickets.session_id = sessions.id").
		Join("movies ON sessions.movie_id = movies.id").
		Where(sq.Eq{
			"tickets.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB).
		QueryRow().
		Scan(&res.ID, &res.User_ID, &res.Price, &res.Session_ID, &res.Title, &res.Seat, &res.Starts_at)

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
		Delete("tickets").
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
func (r *Repository) RetrieveAll(ctx context.Context) ([]internal.Identifiable, error) {

	rows, err := sq.
		Select("tickets.id", "user_id", "price", "session_id", "movies.name", "tickets.seat", "sessions.starts_at").
		From("tickets").
		Join("sessions ON tickets.session_id = sessions.id").
		Join("movies ON sessions.movie_id = movies.id").
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

		err = rows.Scan(&res.ID, &res.User_ID, &res.Price, &res.Session_ID, &res.Title, &res.Seat, &res.Starts_at)
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

func (r *Repository) SeatNumber(id int64, ctx context.Context, tx *sql.Tx) (internal.Identifiable, error) {
	var seat sql.NullInt64

	var res Resource

	err := sq.
		Select("MAX(tickets.seat)").
		From("tickets").
		Join("sessions ON tickets.session_id = sessions.id").
		Join("halls ON sessions.hall_id = halls.id").
		Where(sq.Eq{
			"sessions.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		QueryRow().
		Scan(&seat)

	if err != nil {
		r.Log.Info("Failed to run Retrieve session query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	res.Seat = seat.Int64

	return &res, nil
}

func (r *Repository) HallSeatNumber(id int64, ctx context.Context, tx *sql.Tx) (internal.Identifiable, error) {
	var seat sql.NullInt64

	var res Resource

	err := sq.
		Select("halls.seats").
		From("sessions").
		Join("halls ON sessions.hall_id = halls.id").
		Where(sq.Eq{
			"sessions.id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		QueryRow().
		Scan(&seat)

	if err != nil {
		r.Log.Info("Failed to run Retrieve session query.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	res.Seat = seat.Int64

	return &res, nil
}
