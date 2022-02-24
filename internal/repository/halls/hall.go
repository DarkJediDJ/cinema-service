package hall

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/darkjedidj/cinema-service/internal"
)

type Repository struct {
	DB *sql.DB
}

type Resource struct {
	ID    int  `json:"ID"`
	VIP   bool `json:"VIP"`
	Seats int  `json:"seats"`
}

func (r *Resource) GID() int {
	return r.ID
}

// Create new Hall in DB
func (r *Repository) Create(i internal.Identifiable) (internal.Identifiable, error) {
	var id int64

	hall, _ := i.(*Resource)
	// CHECK ERROR -> hall, ok := i.(*Resource)
	// TODO the above!

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
		// log the actual error here

		return nil, internal.ErrInternalFailure
	}

	return r.Retrieve(id)
}

// Retrieve Hall from DB
func (r *Repository) Retrieve(id int64) (internal.Identifiable, error) {
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
		// TODO logs

		return nil, internal.ErrInternalFailure // TODO do the same EVERYWHERE
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
func (r *Repository) RetrieveAll() ([]internal.Identifiable, error) {
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

		err = rows.Scan(&res.VIP, &res.ID, &res.Seats)
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
