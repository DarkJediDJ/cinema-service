package halls

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/halls"
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
func (s *Service) Create(r internal.Identifiable, ctx context.Context) (internal.Identifiable, error) {
	return s.repo.Create(r, ctx)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(id int64, ctx context.Context) (internal.Identifiable, error) {
	return s.repo.Retrieve(id, ctx)
}

// RetriveAll logic layer for repository method
func (s *Service) RetrieveAll(ctx context.Context) ([]internal.Identifiable, error) {
	return s.repo.RetrieveAll(ctx)
}

// Delete logic layer for repository method
func (s *Service) Delete(id int64, ctx context.Context) error {
	return s.repo.Delete(id, ctx)
}
