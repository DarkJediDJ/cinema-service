package tickets

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/tickets"
)

// Service is a struct to store DB and logger connection
type Service struct {
	repo    *h.Repository
	log     *zap.Logger
	tokenID int64
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
		s.log.Info("Failed to assert ticket object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	if s.tokenID != res.User_ID {
		return nil, internal.ErrValidationFailed
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	tx, err := s.repo.DB.BeginTx(timeoutCtx, nil)
	if err != nil {

		err = tx.Rollback()
		if err != nil {
			return nil, internal.ErrInternalFailure
		}
		return nil, fmt.Errorf("%w:couldn't open transaction connection", err)
	}

	lastSeat, err := s.repo.SeatNumber(res.Session_ID, ctx, tx)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, internal.ErrInternalFailure
		}
		return nil, err
	}

	lSeat, ok := lastSeat.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert ticket object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	maxSeat, err := s.repo.HallSeatNumber(res.Session_ID, ctx, tx)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, internal.ErrInternalFailure
		}
		return nil, err
	}

	mSeat, ok := maxSeat.(*h.Resource)
	if !ok {
		s.log.Info("Failed to assert ticket object.",
			zap.Bool("ok", ok),
		)

		return nil, internal.ErrInternalFailure
	}

	if lSeat.Seat >= mSeat.Seat {
		return nil, internal.ErrNoSeats
	}

	res.Seat = lSeat.Seat + 1

	createdID, err := s.repo.Create(ctx, res, tx)
	if err != nil {
		return nil, internal.ErrInternalFailure
	}

	err = tx.Commit()
	if err != nil {
		return nil, internal.ErrInternalFailure
	}

	return s.repo.Retrieve(int64(createdID), ctx)
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
