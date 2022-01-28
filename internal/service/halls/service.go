package halls

import (
	"database/sql"
	hall "github.com/darkjedidj/cinema-service/internal/repository/halls"
)

type Service struct {
	repo *hall.Repository
}

func Init(db *sql.DB) *Service {
	return &Service{
		repo: &hall.Repository{DB: db},
	}
}
