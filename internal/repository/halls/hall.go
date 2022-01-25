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

func (r *Repository) Create(hall Resource) {
	insertHall := `insert into halls("vip") values($1)`
	_, err := r.DB.Exec(insertHall, hall.VIP)
	if err != nil {
		panic(err)
	}
}

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

func (r *Repository) Delete(hall Resource) {
	insertHall := `DELETE FROM public.halls WHERE "id" = $1`
	_, err := r.DB.Exec(insertHall, hall.ID)
	if err != nil {
		panic(err)
	}
}
