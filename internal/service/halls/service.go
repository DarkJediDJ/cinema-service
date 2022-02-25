package halls

import (
	"database/sql"

	"github.com/darkjedidj/cinema-service/internal"

	h "github.com/darkjedidj/cinema-service/internal/repository/halls"
)

type Service struct {
	repo *h.Repository
}

// Init returns Service object
func Init(db *sql.DB) *Service {
	return &Service{
		repo: &h.Repository{DB: db},
	}
}

// Create logic layer for repository method
func (s *Service) Create(r internal.Identifiable) (internal.Identifiable, error) {
	return s.repo.Create(r)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64) (internal.Identifiable, error) {
	return s.repo.Retrieve(id)
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll() ([]internal.Identifiable, error) {
	return s.repo.RetrieveAll()
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
