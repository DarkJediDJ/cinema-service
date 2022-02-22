package privileges

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
	ID   int    `json:"ID"`
	Name string `json:"Privilage name"`
}

func (r *Resource) GID() int {
	return r.ID
}

// Create new privilage in DB
func (r *Repository) Create(i internal.Identifiable) (internal.Identifiable, error) {
	var id int

	privilage, ok := i.(*Resource)
	if !ok {
		return nil, internal.ErrInternalFailure
	}

	query := sq.
		Insert("privileges").
		Columns("name").
		Values(privilage.Name).
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

// Retrieve privilage from DB
func (r *Repository) Retrieve(id int64) (internal.Identifiable, error) {
	var res Resource

	query := sq.
		Select("id", "name").
		From("privileges").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	err := query.
		QueryRow().
		Scan(&res.ID, &res.Name)

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

// Delete privilage in DB
func (r *Repository) Delete(id int64) error {
	query := sq.
		Delete("privileges").
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

// RetrieveAll privilages from DB
func (r *Repository) RetrieveAll() ([]internal.Identifiable, error) {
	query := sq.
		Select("id","name").
		From("privileges").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.DB)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	var data []*Resource

	for rows.Next() {
		res := &Resource{}

		err = rows.Scan(&res.ID, &res.Name)
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
