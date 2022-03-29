package sessions

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
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
func (s *Service) Create(i internal.Identifiable, ctx context.Context) (internal.Identifiable, error) {
	res, ok := i.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert session object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	valid, err := s.repo.TimeValid(res, ctx)
	if err != nil {
		return nil, err
	}
	if valid {
		return s.repo.Create(i, ctx)
	}
	return nil, fmt.Errorf("%w: this time is already in use", internal.ErrValidationFailed)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64, ctx context.Context) (internal.Identifiable, error) {
	return s.repo.Retrieve(int64(id), ctx)
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll(ctx context.Context) ([]internal.Identifiable, error) {
	return s.repo.RetrieveAll(ctx)
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64, ctx context.Context) error {
	return s.repo.Delete(int64(id), ctx)
}
