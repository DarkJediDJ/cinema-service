package hall

import (
	"database/sql"
	"fmt"
	"github.com/darkjedidj/cinema-service/internal"
	"log"
	"testing"

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

func TestRetrieve(t *testing.T) {
	t.Run("Test retrive", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}
		defer func() {
			repo.DB.Close()
		}()

		query := "SELECT vip, id, seats FROM halls WHERE id = \\$1"
		rows := sqlmock.NewRows([]string{"vip", "id", "seats"}).
			AddRow(hall.VIP, hall.ID, hall.Seats)

		mock.ExpectQuery(query).WithArgs(hall.ID).WillReturnRows(rows)

		halldb, err := repo.Retrieve(int64(hall.ID))

		assert.NotNil(t, halldb)
		assert.NoError(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Test create", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}

		defer func() {
			repo.DB.Close()
		}()

		mock.ExpectQuery("INSERT INTO halls (.*)").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(hall.ID))

		query := "SELECT vip, id, seats FROM halls WHERE id = \\$1"

		rows := sqlmock.NewRows([]string{"vip", "id", "seats"}).
			AddRow(hall.VIP, hall.ID, hall.Seats)

		mock.ExpectQuery(query).WithArgs(hall.ID).WillReturnRows(rows)

		halldb, err := repo.Create(*hall)

		assert.NotNil(t, halldb)
		assert.NoError(t, err)
	})
}

func TestRetrieveErr(t *testing.T) {
	t.Run("Test retrive error", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}

		defer func() {
			repo.DB.Close()
		}()

		query := "SELECT vip, id, seats FROM halls WHERE id = \\$1"
		rows := sqlmock.NewRows([]string{"vip", "id", "seats"})

		mock.ExpectQuery(query).WithArgs(hall.ID).WillReturnRows(rows)

		halldb, err := repo.Retrieve(int64(hall.ID))

		assert.Nil(t, halldb)
		assert.Error(t, err)
	})
}

func TestCreateErr(t *testing.T) {
	t.Run("Test create error", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}

		defer func() {
			repo.DB.Close()
		}()

		mock.ExpectQuery("INSERT INTO halls (.*)").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(hall.ID))

		mock.ExpectQuery("SELECT vip, id, seats FROM halls WHERE id = \\$1").
			WithArgs(hall.ID).
			WillReturnRows(sqlmock.NewRows([]string{"vip", "id", "seats"}))

		halldb, err := repo.Create(*hall)

		assert.Nil(t, halldb)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test delete", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}
		defer func() {
			repo.DB.Close()
		}()

		query := "DELETE FROM halls WHERE id = \\$1"

		mock.ExpectExec(query).WithArgs(hall.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(int64(hall.ID))
		assert.NoError(t, err)
	})
}

func TestDeleteErr(t *testing.T) {
	t.Run("Test delete err", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}
		defer func() {
			repo.DB.Close()
		}()

		query := "DELETE FROM halls WHERE id = \\$1"

		mock.ExpectExec(query).WithArgs(hall.ID)

		err := repo.Delete(int64(hall.ID))
		assert.Error(t, err)
	})
}

func TestRetrieveNil(t *testing.T) {
	t.Run("Test retrive no rows", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}
		defer func() {
			repo.DB.Close()
		}()

		query := "SELECT vip, id, seats FROM halls WHERE id = \\$1"
		rows := sqlmock.NewRows([]string{"vip", "id", "seats"}).
			AddRow(hall.VIP, hall.ID, hall.Seats)

		mock.ExpectQuery(query).WithArgs(hall.ID).WillReturnRows(rows)

		halldb, err := repo.Retrieve(int64(1))
		assert.Nil(t, halldb)
		assert.Error(t, err)
	})
}

func TestRetrieveAll(t *testing.T) {
	t.Run("Test retriveAll", func(t *testing.T) {
		db, mock := NewMock()
		repo := &Repository{DB: db}
		defer func() {
			repo.DB.Close()
		}()

		query := "SELECT vip, id, seats FROM halls"
		rows := sqlmock.NewRows([]string{"vip", "id", "seats"}).
			AddRow(hall.VIP, hall.ID, hall.Seats)

		mock.ExpectQuery(query).WillReturnRows(rows)

		halldb, err := repo.RetrieveAll()
		assert.NotNil(t, halldb)
		assert.NoError(t, err)
	})
}

func TestCreateNil(t *testing.T) {
	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	testCreateNilCases := []struct {
		name           string
		expectedError  error
		expectedResult *Resource
		prepare        func(sqlm2 sqlmock.Sqlmock)
	}{
		{
			name:           "failed",
			expectedError:  internal.ErrInternalFailure,
			expectedResult: nil,
			prepare: func(sqlm2 sqlmock.Sqlmock) {
				sqlm2.ExpectQuery("INSERT INTO halls (.*)").
					WillReturnError(fmt.Errorf("blah blah"))
			},
		},
	}

	for _, tc := range testCreateNilCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &Repository{DB: db}

			tc.prepare(mock)
			res, err := repo.Create(*hall)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
