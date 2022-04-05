package user

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

var user = &Resource{
	ID:       15,
	EMail:    "mail@gmail.com",
	Password: "password",
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
				sqlm2.ExpectExec("INSERT INTO users (email,password) VALUES ($1,$2)").
					WillReturnError(internal.ErrInternalFailure)
			},
			object: user,
		},
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: user,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email,password) VALUES ($1,$2)")).WithArgs(user.EMail, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			object: user,
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("INSERT INTO users (email,password) VALUES ($1,$2)").
					WillReturnError(fmt.Errorf("unable to retrieve Resource"))
			},
			object: user,
		},
		{
			name:           "failed, assertion error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec(regexp.QuoteMeta("INSERT INTO users (email,password) VALUES ($1,$2)")).
					WillReturnError(internal.ErrInternalFailure)
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

			tc.prepare(mock)
			err = repo.Create(tc.object, ctx)

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
			expectedResult: user,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT id, email, password FROM users WHERE email = \\$1").
					WithArgs(user.EMail).
					WillReturnRows(sqlm2.
						NewRows([]string{"id", "email", "password"}).
						AddRow(user.ID, user.EMail, user.Password))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT id, email, password FROM users WHERE email = \\$1").
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT id, email, password FROM users WHERE email = \\$1").
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
			res, err := repo.Retrieve(user.EMail, ctx)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRetrievePrivileges(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testRetrieveCases := []struct {
		name           string
		expectedError  error
		expectedResult []string
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: []string{"hall"},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT privileges.name FROM privileges JOIN user_privileges on user_privileges.privilege_id = privileges.id WHERE user_privileges.user_id = \\$1").
					WithArgs(user.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"privileges.name"}).
						AddRow("hall"))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT privileges.name FROM privileges JOIN user_privileges on user_privileges.privilege_id = privileges.id WHERE user_privileges.user_id = \\$1").
					WithArgs(user.ID).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
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

			tc.prepare(mock)
			res, err := repo.RetrievePrivileges(user.ID)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRetrieveTickets(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testRetrieveCases := []struct {
		name           string
		expectedError  error
		expectedResult bool
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: true,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id FROM tickets WHERE tickets.id = \\$1 AND tickets.user_id = \\$2").
					WithArgs(1, user.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"tickets.id"}).
						AddRow(1))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: false,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT tickets.id FROM tickets WHERE tickets.id = \\$1 AND tickets.user_id = \\$2").
					WithArgs(1, user.ID).
					WillReturnError(fmt.Errorf("unable to perform your request, please try again later"))
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

			tc.prepare(mock)
			res, err := repo.RetrieveTickets(1, user.ID)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGID(t *testing.T) {
	res := &Resource{ID: user.ID}
	assert.Equal(t, user.ID, res.GID())
}
