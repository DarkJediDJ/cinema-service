package user_privileges

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

var user_privileges = &Resource{
	ID:           1,
	User_id:      0,
	Privilege_id: 0,
	Email:        "0",
	Privilege:    "privilege",
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
				sqlm2.ExpectQuery("INSERT INTO user_privileges (.*)").
					WillReturnError(internal.ErrInternalFailure)
			},
			object: user_privileges,
		},
		{
			name:           "success",
			expectedError:  nil,
			expectedResult: user_privileges,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO user_privileges (.*)").
					WillReturnRows(sqlm2.
						NewRows([]string{"id"}).
						AddRow(user_privileges.ID))
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id WHERE user_privileges.id = \\$1").
					WithArgs(user_privileges.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"user.email", "id", "privileges.name"}).
						AddRow(user_privileges.Email, user_privileges.ID, user_privileges.Privilege))
			},
			object: user_privileges,
		},
		{
			name:           "failed, retrieve error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO user_privileges (.*)").
					WillReturnError(fmt.Errorf("unable to retrieve Resource"))
			},
			object: user_privileges,
		},
		{
			name:           "failed, assertion error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO user_privileges (.*)").
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
			expectedResult: user_privileges,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id WHERE user_privileges.id = \\$1").
					WithArgs(user_privileges.ID).
					WillReturnRows(sqlm2.
						NewRows([]string{"user.email", "id", "privileges.name"}).
						AddRow(user_privileges.Email, user_privileges.ID, user_privileges.Privilege))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id WHERE user_privileges.id = \\$1").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id WHERE user_privileges.id = \\$1").
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
			res, err := repo.Retrieve(user_privileges.ID, ctx)

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
			expectedResult: []internal.Identifiable{user_privileges},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id").
					WillReturnRows(sqlm2.
						NewRows([]string{"user.email", "id", "privileges.name"}).
						AddRow(user_privileges.Email, user_privileges.ID, user_privileges.Privilege))
			},
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id").
					WillReturnError(internal.ErrInternalFailure)
			},
		},
		{
			name:           "failed, sql no rows error",
			expectedError:  nil,
			expectedResult: []internal.Identifiable{},
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("SELECT users.email, user_privileges.id, privileges.name FROM user_privileges JOIN users ON user_privileges.user_id = users.id JOIN privileges ON user_privileges.privilege_id = privileges.id").
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
			expectedResult: user_privileges,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM user_privileges WHERE id = \\$1").
					WithArgs(user_privileges.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: user_privileges.ID,
		},
		{
			name:           "failed, database error",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectExec("DELETE FROM user_privileges WHERE id = \\$1").
					WithArgs(user_privileges.ID).
					WillReturnError(internal.ErrInternalFailure)
			},
			id: int64(user_privileges.ID),
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
	res := &Resource{ID: user_privileges.ID}
	assert.Equal(t, user_privileges.ID, res.GID())
}
