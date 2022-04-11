package users

import (
	"context"
	"database/sql"
	"regexp"

	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/users"
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
func (s *Service) Create(i internal.Identifiable, ctx context.Context) error {
	res, ok := i.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert ticket object.",
			zap.Bool("ok", ok),
		)

		return internal.ErrInternalFailure
	}

	match, err := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, res.EMail)
	if err != nil {
		return internal.ErrInternalFailure
	}
	if !match {
		return internal.ErrWrongEmail
	}

	dbuser, err := s.repo.Retrieve(res.EMail, ctx)
	if err != nil {
		return internal.ErrInternalFailure
	}

	if dbuser != nil {
		return internal.ErrValidationFailed
	}

	return s.repo.Create(res, ctx)
}

// Retrieve logic layer for repository method
func (s *Service) Retrieve(email string, ctx context.Context) (internal.Identifiable, error) {
	return s.repo.Retrieve(email, ctx)
}

// RetrievePrivileges logic layer for repository method
func (s *Service) RetrievePrivileges(id int64) ([]string, error) {
	return s.repo.RetrievePrivileges(id)
}

// RetrieveTickets logic layer for repository method
func (s *Service) RetrieveTickets(ticket int64, user int64) (bool, error) {
	return s.repo.RetrieveTickets(ticket, user)
}
