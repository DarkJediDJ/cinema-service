package hall

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
)

var hall = &Resource{
	ID:    15,
	VIP:   true,
	Seats: 15,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testCreateCases := []struct {
		name           string
		expectedError  error
		expectedResult internal.Identifiable
		prepare        func(sqlm2 sqlmock.Sqlmock)
		object         internal.Identifiable
	}{
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
			object: hall,
		},
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: hall,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(hall.ID))
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WithArgs(hall.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"vip", "id", "seats"}).
						AddRow(hall.VIP, hall.ID, hall.Seats))
			},
			object: hall,
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnError(fmt.Errorf("unable to retrieve Resource"))
			},
			object: hall,
		},
		{
			name:           "failed, assertion error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
			object: nil,
		},
	}

	for _, tc := range testCreateCases {
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
			res, err := repo.Create(tc.object, ctx)

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
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: hall,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WithArgs(hall.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"vip", "id", "seats"}).
						AddRow(hall.VIP, hall.ID, hall.Seats))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WillReturnRows(sqlm2.
						NewRows(nil))
			},
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
			res, err := repo.Retrieve(hall.ID, ctx)

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
			expectedResult: []internal.Identifiable{hall},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls").
					WillReturnRows(sqlm2.
						NewRows([]string{"vip", "id", "seats"}).
						AddRow(hall.VIP, hall.ID, hall.Seats))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: []internal.Identifiable{},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls").
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
			expectedResult: hall,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM halls WHERE id = \\$1").
					WithArgs(hall.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: int64(hall.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM halls WHERE id = \\$1").
					WithArgs(hall.ID).
					WillReturnError(internal.ErrInternalFailure)
			},
			id: int64(hall.ID),
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

func TestGID(t *testing.T) {
	res := &Resource{ID: hall.ID}
	assert.Equal(t, hall.ID, res.GID())
}
