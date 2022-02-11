package movie

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	DB *sql.DB
}

type Resource struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Duration string `json:"Duration"`
}

// Create new Movie in DB
func (r *Repository) Create(movie Resource) (*Resource, error) {
	var id int

	query := sq.
		Insert("movies").
		Columns("name", "duration").
		Values(movie.Name, movie.Duration).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err := query.
		QueryRow().
		Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.Retrieve(int64(id))
}

// Retrieve Movie from DB
func (r *Repository) Retrieve(id int64) (*Resource, error) {
	var res Resource

	query := sq.
		Select("name", "duration", "id").
		From("movies").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err := query.
		QueryRow().
		Scan(&res.Name, &res.Duration, &res.ID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Delete Movie in DB
func (r *Repository) Delete(id int64) error {
	query := sq.
		Delete("movies").
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

// RetrieveAll Movies from DB
func (r *Repository) RetrieveAll() ([]*Resource, error) {
	query := sq.
		Select("name", "duration", "id").
		From("movies").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(&res.Name, &res.Duration, &res.ID)

		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		data = append(data, res)
	}

	return data, nil
}
