package halls

import (
	"database/sql"

	"github.com/darkjedidj/cinema-service/internal"
	"go.uber.org/zap"

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
