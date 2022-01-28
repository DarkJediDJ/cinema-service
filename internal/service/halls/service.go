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

func (s *Service) Create(hall hall.Resource) error {
	err := s.repo.Create(hall)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Retrieve(id int64) (hall.Resource, error) {
	hall, err := s.repo.Retrieve(int64(id))
	if err != nil {
		return hall, err
	}
	return hall, nil
}

func (s *Service) RetrieveAll() ([]hall.Resource, error) {
	halls, err := s.repo.RetrieveAll()
	if err != nil {
		return nil, err
	}
	return halls, nil
}

func (s *Service) Delete(id int64) error {
	err := s.repo.Delete(int64(id))
	if err != nil {
		return err
	}
	return nil
}
