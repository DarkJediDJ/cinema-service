package hall

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/darkjedidj/cinema-service/internal"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
	}{
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
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
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(hall.ID)).WillReturnError(fmt.Errorf("unable to retrieve Resource"))
			},
		},
	}

	for _, tc := range testCreateCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.Create(hall)

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
			expectedResult: hall,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WithArgs(hall.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"vip", "id", "seats"}).
						AddRow(hall.VIP, hall.ID, hall.Seats))
			},
			id: int64(hall.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(hall.ID),
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
					WillReturnRows(sqlm2.
						NewRows([]string{}))
			},
			id: 5,
		},
	}

	for _, tc := range testRetrieveCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.Retrieve(tc.id)

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
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
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
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.RetrieveAll()
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
		expectedResult *Resource
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
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(hall.ID),
		},
	}

	for _, tc := range testDeleteCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			err := repo.Delete(tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
