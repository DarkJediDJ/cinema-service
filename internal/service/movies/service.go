package movies

import (
	"database/sql"
	"time"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/movies"
	"go.uber.org/zap"
)

const maxMinutes, minMinutes, maxLetters, minLetters = 350, 30, 0, 50

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
func (s *Service) Create(r internal.Identifiable) (internal.Identifiable, error) {
	entity, ok := r.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert movie object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	duration, err := time.ParseDuration(entity.Duration)

	if err != nil {
		s.log.Info("Failed to parse time.",
			zap.Error(err),
		)

		return nil, internal.ErrInternalFailure
	}

	if duration.Minutes() < minMinutes || duration.Minutes() > maxMinutes {
		s.log.Info("Wrond duration value")

		return nil, internal.ErrInternalFailure
	}

	if len(entity.Name) <= minLetters || len(entity.Name) > maxLetters {
		s.log.Info("Wrong name value")

		return nil, internal.ErrInternalFailure
	}

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
