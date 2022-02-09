package movies

import (
	"database/sql"
	"errors"
	"time"

	h "github.com/darkjedidj/cinema-service/internal/repository/movies"
)

type Service struct {
	repo *h.Repository
}

const minMinutes = 30
const maxMinutes = 350
const minLetters = 0
const maxLetters = 50

// Init returns Service object
func Init(db *sql.DB) *Service {
	return &Service{
		repo: &h.Repository{DB: db},
	}
}

// Create logic layer for repository method
func (s *Service) Create(movie h.Resource) (*h.Resource, error) {
	duration, err := time.ParseDuration(movie.Duration)

	if err != nil {
		return nil, err
	}

	if duration.Minutes() < minMinutes || duration.Minutes() > maxMinutes {
		return nil, errors.New("duration of movie is incorrect")
	}

	if len(movie.Name) <= minLetters || len(movie.Name) > maxLetters {
		return nil, errors.New("name of movie is incorrect")
	}

	return s.repo.Create(movie)
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
