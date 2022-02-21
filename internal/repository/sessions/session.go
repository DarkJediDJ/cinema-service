package session

import (
	"database/sql"
	"fmt"

	"github.com/darkjedidj/cinema-service/internal"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	DB *sql.DB
}

type Resource struct {
	ID       int    `json:"ID"`
	Hall_id  int    `json:"hall_id,omitempty"`
	Movie_id int    `json:"movie_id,omitempty"`
	Schedule string `json:"Schedule"`
	VIP      bool   `json:"VIP"`
	Name     string `json:"Movie name"`
}

func (r *Resource) GID() int {
	return r.ID
}

// Create new session in DB
func (r *Repository) Create(i internal.Identifiable) (internal.Identifiable, error) {
	var id int

	session, ok := i.(*Resource)
	if !ok {
		return nil, internal.ErrInternalFailure
	}

	query := sq.
		Insert("sessions").
		Columns("hall_id", "movie_id", "schedule").
		Values(session.Hall_id, session.Movie_id, session.Schedule).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)
	err := query.
		QueryRow().
		Scan(&id)

	if err != nil {
		fmt.Print(err)
		return nil, internal.ErrInternalFailure
	}

	return r.Retrieve(int64(id))
}

// Retrieve session from DB
func (r *Repository) Retrieve(id int64) (internal.Identifiable, error) {
	var res Resource

	query := sq.
		Select("sessions.id", "halls.vip", "movies.name", "schedule").
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
		Scan(&res.ID, &res.VIP, &res.Name, &res.Schedule)

	if err == sql.ErrNoRows {
		fmt.Print(err)
		return nil, nil
	}

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return &res, nil
}

// Delete session in DB
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
		return err
	}

	return nil
}

// RetrieveAll sessions from DB
func (r *Repository) RetrieveAll() ([]internal.Identifiable, error) {
	query := sq.
		Select("sessions.id", "halls.vip", "movies.name", "schedule").
		From("sessions").
		Join("movies ON sessions.movie_id = movies.id").
		Join("halls ON sessions.hall_id = halls.id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(&res.ID, &res.VIP, &res.Name, &res.Schedule)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
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
