package hall

import (
	"database/sql"
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
func (r *Repository) Create(hall Resource) error {
	insertHall := `insert into halls("vip","seats") values($1,$2)`
	_, err := r.DB.Exec(insertHall, hall.VIP, hall.Seats)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve Hall from DB
func (r *Repository) Retrieve(id int64) (dbHall Resource, e error) {

	rows, err := r.DB.Query(`SELECT * FROM halls WHERE "id" = $1`, id)
	if err != nil {
		e = err
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dbHall.VIP, &dbHall.ID, &dbHall.Seats)
		if err != nil {
			e = err
			return
		}
	}
	e = nil
	return
}

// Delete Hall in DB
func (r *Repository) Delete(id int64) error {
	insertHall := `DELETE FROM public.halls WHERE "id" = $1`
	_, err := r.DB.Exec(insertHall, id)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveAll halls from DB
func (r *Repository) RetrieveAll() ([]Resource, error) {
	var hall Resource
	sqlStatement := `SELECT vip, id, seats
   FROM public.halls`
	rows, err := r.DB.Query(sqlStatement)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	var hallSlice []Resource
	for rows.Next() {
		err = rows.Scan(&hall.VIP, &hall.ID, &hall.Seats)
		if err != nil {
			return nil,err
		}
		hallSlice = append(hallSlice, hall)
	}
	return hallSlice,nil
}
