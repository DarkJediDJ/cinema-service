package movies

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/movies"
)

const maxMinutes, minMinutes, maxLetters, minLetters = 350, 30, 50, 0

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

	duration, err := time.ParseDuration(res.Duration)
	if err != nil {
		error := fmt.Errorf("%w: failed to parse duration", internal.ErrValidationFailed)

		return nil, error
	}

	if duration.Minutes() < minMinutes {
		error := fmt.Errorf("%w: duration too short", internal.ErrValidationFailed)

		return nil, error
	}

	if duration.Minutes() > maxMinutes {
		error := fmt.Errorf("%w: duration too long", internal.ErrValidationFailed)

		return nil, error
	}

	if len(res.Name) <= minLetters {
		error := fmt.Errorf("%w: name too short", internal.ErrValidationFailed)

		return nil, error
	}

	if len(res.Name) > maxLetters {
		error := fmt.Errorf("%w: name too long", internal.ErrValidationFailed)

		return nil, error
	}

	return s.repo.Create(i)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64) (internal.Identifiable, error) {
	return s.repo.Retrieve(id)
}

// RetrieveAll logic layer for repository method
func (s *Service) RetrieveAll() ([]internal.Identifiable, error) {
	return s.repo.RetrieveAll()
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
