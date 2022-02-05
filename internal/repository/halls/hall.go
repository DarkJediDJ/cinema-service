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
func (r *Repository) Create(hall Resource) (*Resource, error) {
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
func (r *Repository) Retrieve(id int64) (*Resource, error) {
	var res Resource

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
		Scan(&res.VIP, &res.ID, &res.Seats)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
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
func (r *Repository) RetrieveAll() ([]*Resource, error) {
	query := sq.
		Select("vip", "id", "seats").
		From("halls").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(res.VIP, res.ID, res.Seats)
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
