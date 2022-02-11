package movie

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/darkjedidj/cinema-service/internal"
	"github.com/stretchr/testify/assert"
)

var movie = &Resource{
	ID:       15,
	Name:     "Lord of the Rings",
	Duration: "2h22m",
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
		expectedResult *Resource
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO movies (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: movie,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO movies (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(movie.ID))
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies WHERE id = \\$1").
					WithArgs(movie.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"name", "duration", "id"}).
						AddRow(movie.Name, movie.Duration, movie.ID))
			},
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO movies (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(movie.ID)).WillReturnError(internal.ErrInternalFailure)
			},
		},
	}

	for _, tc := range testCreateCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.Create(*movie)

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
		expectedResult *Resource
		prepare        func(sqlm2 sqlmock.Sqlmock)
		id             int64
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: movie,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies WHERE id = \\$1").
					WithArgs(movie.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"name", "duration", "id"}).
						AddRow(movie.Name, movie.Duration, movie.ID))
			},
			id: int64(movie.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies WHERE id = \\$1").
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(movie.ID),
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies WHERE id = \\$1").
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
		expectedResult []*Resource
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: []*Resource{movie},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies").
					WillReturnRows(sqlm2.
						NewRows([]string{"name", "duration", "id"}).
						AddRow(movie.Name, movie.Duration, movie.ID))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT name, duration, id FROM movies").
					WillReturnRows(sqlm2.NewRows([]string{}))
			},
		},
	}

	for _, tc := range testRetrieveAllCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.RetrieveAll()
			fmt.Print(res)
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
			expectedResult: movie,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM movies WHERE id = \\$1").
					WithArgs(movie.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: int64(movie.ID),
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM movies WHERE id = \\$1").
					WithArgs(movie.ID).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
			id: int64(movie.ID),
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
