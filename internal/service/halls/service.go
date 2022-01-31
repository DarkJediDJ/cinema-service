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
func (s *Service) Create(hall h.Resource) (h.Resource, error) {
	dbhall, err := s.repo.Create(hall)
	if err != nil {
		return h.Resource{}, err
	}

	return dbhall, nil
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64) (h.Resource, error) {
	dbhall, err := s.repo.Retrieve(int64(id))
	if err != nil {
		return h.Resource{}, err
	}

	return dbhall, nil
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll() ([]h.Resource, error) {
	halls, err := s.repo.RetrieveAll()
	if err != nil {
		return nil, err
	}

	return halls, nil
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64) error {
	err := s.repo.Delete(int64(id))
	if err != nil {
		return err
	}

	return nil
}
