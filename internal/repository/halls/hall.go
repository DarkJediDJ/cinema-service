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
func (r *Repository) Create(hall Resource) (dbHall Resource, e error) {
	create := sq.
		Insert("halls").
		Columns("vip", "seats").
		Values(hall.VIP, hall.Seats).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	_, err := create.
		Exec()

	if err != nil {
		return Resource{}, err
	}

	retrive := sq.
		Select("vip", "id", "seats").
		From("halls").
		OrderBy("id DESC").Limit(1).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err = retrive.
		QueryRow().
		Scan(&dbHall.VIP, &dbHall.ID, &dbHall.Seats)
	
	if err != nil {
		return Resource{}, err
	}

	e = nil
	return
}

// Retrieve Hall from DB
func (r *Repository) Retrieve(id int64) (dbHall Resource, e error) {
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

	if err != nil {
		return Resource{}, err
	}

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
		if err != nil {
			return nil, err
		}

		hallSlice = append(hallSlice, hall)
	}

	return hallSlice, nil
}
