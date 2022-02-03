package hall

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	DB *sql.DB
}

type Resource struct {
	ID    int  `json:"ID"`
	VIP   bool `json:"VIP"`
	Seats int  `json:"seats"`
}

// Create new Hall in DB
func (r *Repository) Create(hall Resource) (dbHall *Resource, e error) {
	var id int

	query := sq.
		Insert("halls").
		Columns("vip", "seats").
		Values(hall.VIP, hall.Seats).
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

// Retrieve Hall from DB
func (r *Repository) Retrieve(id int64) (hall *Resource, e error) {
	var dbHall Resource
	query := sq.
		Select("vip", "id", "seats").
		From("halls").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err := query.
		QueryRow().
		Scan(&dbHall.VIP, &dbHall.ID, &dbHall.Seats)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	hall = &dbHall
	e = nil
	return
}

// Delete Hall in DB
func (r *Repository) Delete(id int64) error {
	query := sq.
		Delete("halls").
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

// RetrieveAll halls from DB
func (r *Repository) RetrieveAll() ([]Resource, error) {
	var hall Resource

	var hallSlice []Resource

	query := sq.
		Select("vip", "id", "seats").
		From("halls").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&hall.VIP, &hall.ID, &hall.Seats)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		hallSlice = append(hallSlice, hall)
	}

	return hallSlice, nil
}
