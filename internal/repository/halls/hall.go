package halldb

import (
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

type Resource struct {
	ID  int  `json:"ID"`
	VIP bool `json:"VIP"`
}

//Create new Hall in DB
func (r *Repository) Create(hall Resource) {
	insertHall := `insert into halls("vip") values($1)`
	_, err := r.DB.Exec(insertHall, hall.VIP)
	if err != nil {
		panic(err)
	}
}

//Retrieve Hall from DB
func (r *Repository) Retrieve(id int64) (dbHall Resource) {

	rows, err := r.DB.Query(`SELECT * FROM halls WHERE "id" = $1`, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&dbHall.VIP, &dbHall.ID); err != nil {
			panic(err)
		}
	}
	return
}

//Delete Hall in DB
func (r *Repository) Delete(id int64) {
	insertHall := `DELETE FROM public.halls WHERE "id" = $1`
	_, err := r.DB.Exec(insertHall, id)
	if err != nil {
		panic(err)
	}
}
