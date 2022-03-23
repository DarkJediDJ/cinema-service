package sessions

import (
	"database/sql"
	"fmt"

	"github.com/darkjedidj/cinema-service/internal"
	"go.uber.org/zap"

	h "github.com/darkjedidj/cinema-service/internal/repository/sessions"
)

// Service is a struct to store DB and logger connection
type Service struct {
	repo *h.Repository
	log  *zap.Logger
}

// Init returns Service object
func Init(db *sql.DB, l *zap.Logger) *Service {

	return &Service{
		repo: &h.Repository{DB: db, Log: l},
		log:  l,
	}
}

// Create logic layer for repository method
func (s *Service) Create(i internal.Identifiable) (internal.Identifiable, error) {
	res, ok := i.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert movie object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	valid, err := s.repo.TimeValid(res)
	if err != nil {
		return nil, err
	}
	if valid {
		return s.repo.Create(i)
	}
	return nil, fmt.Errorf("%w: this time is already in use", internal.ErrValidationFailed)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64) (internal.Identifiable, error) {
	return s.repo.Retrieve(int64(id))
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll() ([]internal.Identifiable, error) {
	return s.repo.RetrieveAll()
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(int64(id))
}
