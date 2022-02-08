package halls

import (
	"database/sql"

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
func (s *Service) Create(hall h.Resource) (*h.Resource, error) {
	return s.repo.Create(hall)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64) (*h.Resource, error) {
	return s.repo.Retrieve(int64(id))
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll() ([]*h.Resource, error) {
	return s.repo.RetrieveAll()
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(int64(id))
}
