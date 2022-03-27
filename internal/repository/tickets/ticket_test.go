package tickets

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
)

var ticket = &Resource{
	Starts_at:  "13:25",
	Price:      12.2,
	Seat:       1,
	ID:         1,
	Title:      "Matrix",
	User_ID:    1,
	Session_ID: 1,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestCreate(t *testing.T) {
	testCreateCases := []struct {
		name              string
		expectedError     error
		expectedResult    int64
		prepare           func(sqlm2 sqlmock.Sqlmock)
		transactionResult func(sqlm2 sqlmock.Sqlmock)
		object            internal.Identifiable
	}{
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: 0,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO tickets (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectRollback()
			},
			object: ticket,
		},
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: 1,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO tickets (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(ticket.ID))
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id WHERE tickets.id = \\$1").
					WithArgs(ticket.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"id", "user_id", " price", "session_id", "name", "seat", "starts_at"}).
						AddRow(ticket.ID, ticket.User_ID, ticket.Price, ticket.Session_ID, ticket.Title, ticket.Seat, ticket.Starts_at))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectCommit()
			},
			object: ticket,
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: 0,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO tickets (.*)").
					WillReturnError(fmt.Errorf("unable to retrieve Resource"))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectRollback()
			},
			object: ticket,
		},
		{
			name:           "failed, assertion error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: 0,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO tickets (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectRollback()
			},
			object: nil,
		},
	}

	for _, tc := range testCreateCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := NewMock()
			defer func() {
				db.Close()
			}()

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			mock.ExpectBegin()
			tx, err := repo.DB.Begin()
			if err != nil {
				log.Fatalf("can't start transaction : %v", err)
			}

			tc.prepare(mock)

			res, err := repo.Create(ctx, tc.object, tx)

			tc.transactionResult(mock)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRetrieve(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testRetrieveCases := []struct {
		name           string
		expectedError  error
		expectedResult internal.Identifiable
		prepare        func(sqlm2 sqlmock.Sqlmock)
		id             int64
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: ticket,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id WHERE tickets.id = \\$1").
					WithArgs(ticket.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"id", "user_id", " price", "session_id", "name", "seat", "starts_at"}).
						AddRow(ticket.ID, ticket.User_ID, ticket.Price, ticket.Session_ID, ticket.Title, ticket.Seat, ticket.Starts_at))
			},
			id: int64(ticket.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id WHERE tickets.id = \\$1").
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(ticket.ID),
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id WHERE tickets.id = \\$1").
					WillReturnRows(sqlm2.
						NewRows(nil))
			},
			id: 5,
		},
	}

	for _, tc := range testRetrieveCases {
		t.Run(tc.name, func(t *testing.T) {

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			tc.prepare(mock)
			res, err := repo.Retrieve(ticket.ID, ctx)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRetrieveAll(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testRetrieveAllCases := []struct {
		name           string
		expectedError  error
		expectedResult []internal.Identifiable
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: []internal.Identifiable{ticket},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id").
					WillReturnRows(sqlm2.
						NewRows([]string{"id", "user_id", " price", "session_id", "name", "seat", "starts_at"}).
						AddRow(ticket.ID, ticket.User_ID, ticket.Price, ticket.Session_ID, ticket.Title, ticket.Seat, ticket.Starts_at))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: []internal.Identifiable{},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id, user_id, price, session_id, movies.name, tickets.seat, sessions.starts_at FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN movies ON sessions.movie_id = movies.id").
					WillReturnRows(sqlm2.NewRows([]string{}))
			},
		},
	}

	for _, tc := range testRetrieveAllCases {
		t.Run(tc.name, func(t *testing.T) {

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			tc.prepare(mock)
			res, err := repo.RetrieveAll(ctx)
			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDelete(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testDeleteCases := []struct {
		name           string
		expectedError  error
		expectedResult internal.Identifiable
		prepare        func(sqlm2 sqlmock.Sqlmock)
		id             int64
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: ticket,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM tickets WHERE id = \\$1").
					WithArgs(ticket.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: int64(ticket.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM tickets WHERE id = \\$1").
					WithArgs(ticket.ID).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(ticket.ID),
		},
	}

	for _, tc := range testDeleteCases {
		t.Run(tc.name, func(t *testing.T) {

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			tc.prepare(mock)
			err = repo.Delete(tc.id, ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSeatNumber(t *testing.T) {
	testSeatNumberCases := []struct {
		name              string
		expectedError     error
		expectedResult    internal.Identifiable
		prepare           func(sqlm2 sqlmock.Sqlmock)
		transactionResult func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:          "success",
			expectedError: nil,
			expectedResult: &Resource{
				Starts_at:  "",
				Price:      0,
				Seat:       1,
				ID:         0,
				Title:      "",
				User_ID:    0,
				Session_ID: 0,
			},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery(regexp.QuoteMeta("SELECT MAX(tickets.seat) FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN halls ON sessions.hall_id = halls.id WHERE sessions.id = $1")).
					WithArgs(ticket.Session_ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"seat"}).
						AddRow(ticket.Seat))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectCommit()
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery(regexp.QuoteMeta("SELECT MAX(tickets.seat) FROM tickets JOIN sessions ON tickets.session_id = sessions.id JOIN halls ON sessions.hall_id = halls.id WHERE sessions.id = $1")).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectRollback()
			},
		},
	}

	for _, tc := range testSeatNumberCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := NewMock()
			defer func() {
				db.Close()
			}()

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			mock.ExpectBegin()
			tx, err := repo.DB.Begin()
			if err != nil {
				log.Fatalf("can't start transaction : %v", err)
			}

			tc.prepare(mock)

			lastSeat, err := repo.SeatNumber(ticket.Session_ID, ctx, tx)

			tc.transactionResult(mock)

			assert.Equal(t, tc.expectedResult, lastSeat)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestHallSeatNumber(t *testing.T) {
	testHallSeatNumberCases := []struct {
		name              string
		expectedError     error
		expectedResult    internal.Identifiable
		prepare           func(sqlm2 sqlmock.Sqlmock)
		transactionResult func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:          "success",
			expectedError: nil,
			expectedResult: &Resource{
				Starts_at:  "",
				Price:      0,
				Seat:       1,
				ID:         0,
				Title:      "",
				User_ID:    0,
				Session_ID: 0,
			},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery(regexp.QuoteMeta("SELECT halls.seats FROM sessions JOIN halls ON sessions.hall_id = halls.id WHERE sessions.id = $1")).
					WithArgs(ticket.Session_ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"seat"}).
						AddRow(ticket.Seat))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectCommit()
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery(regexp.QuoteMeta("SELECT halls.seats FROM sessions JOIN halls ON sessions.hall_id = halls.id WHERE sessions.id = $1")).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			transactionResult: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectRollback()
			},
		},
	}

	for _, tc := range testHallSeatNumberCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := NewMock()
			defer func() {
				db.Close()
			}()

			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %v", err)
			}

			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println(err)
				}
			}()

			repo := &Repository{DB: db, Log: logger}
			ctx := context.Background()

			mock.ExpectBegin()
			tx, err := repo.DB.Begin()
			if err != nil {
				log.Fatalf("can't start transaction : %v", err)
			}

			tc.prepare(mock)

			lastSeat, err := repo.HallSeatNumber(ticket.Session_ID, ctx, tx)

			tc.transactionResult(mock)

			assert.Equal(t, tc.expectedResult, lastSeat)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGID(t *testing.T) {
	res := &Resource{ID: ticket.ID}
	assert.Equal(t, ticket.ID, res.GID())
}
