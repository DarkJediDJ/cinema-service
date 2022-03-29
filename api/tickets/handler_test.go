package tickets

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	movie "github.com/darkjedidj/cinema-service/internal/repository/tickets"
	"github.com/darkjedidj/cinema-service/test"
)

func TestCreate(t *testing.T) {
	testCreateCases := []struct {
		name           string
		mockService    *test.MockService
		body           string
		id             int64
		expectedStatus int
	}{
		{
			name: "failure: empty body",
			mockService: &test.MockService{
				ExpectedResult: &movie.Resource{
					Starts_at:  "13:25",
					Price:      12.2,
					Seat:       1,
					ID:         1,
					Title:      "Matrix",
					User_ID:    1,
					Session_ID: 1,
				},
			},
			id:             4,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedResult: &movie.Resource{
					Starts_at:  "13:25",
					Price:      12.2,
					Seat:       1,
					ID:         1,
					Title:      "Matrix",
					User_ID:    1,
					Session_ID: 1,
				},
			},
			body: `{
				"User_i":  4,
				"Price": 12.2
			},`,
			id:             4,
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			body: `{
				"User_i":  4,
				"Price": 12.2
			},`,
			id:             4,
			expectedStatus: http.StatusUnprocessableEntity,
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

			w := httptest.NewRecorder()

			vars := map[string]string{
				"id": fmt.Sprintf("%d", tc.id),
			}

			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8085/v1/sessions/%d/tickets", tc.id), strings.NewReader(tc.body))

			r = mux.SetURLVars(r, vars)

			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService, log: logger}).Create(w, r)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}

func TestRetrieve(t *testing.T) {
	testRetrieveCases := []struct {
		name           string
		mockService    *test.MockService
		id             int64
		expectedStatus int
	}{
		{
			name: "failure: no rows",
			mockService: &test.MockService{
				ExpectedResult: nil,
			},
			id:             20,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedResult: &movie.Resource{
					Starts_at:  "13:25",
					Price:      12.2,
					Seat:       1,
					ID:         1,
					Title:      "Matrix",
					User_ID:    1,
					Session_ID: 1,
				},
			},
			id:             15,
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			id:             15,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testRetrieveCases {

		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}

		defer func() {
			if err := logger.Sync(); err != nil {
				fmt.Println(err)
			}
		}()

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			vars := map[string]string{
				"id": fmt.Sprintf("%d", tc.id),
			}

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%d", tc.id), nil)

			r = mux.SetURLVars(r, vars)

			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService}).HandleID(w, r)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}

func TestRetrieveAll(t *testing.T) {
	testRetrieveAllCases := []struct {
		name           string
		mockService    *test.MockService
		expectedStatus int
	}{
		{
			name: "failure: no rows",
			mockService: &test.MockService{
				ExpectedArray: nil,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedArray: []internal.Identifiable{&movie.Resource{
					Starts_at:  "13:25",
					Price:      12.2,
					Seat:       1,
					ID:         1,
					Title:      "Matrix",
					User_ID:    1,
					Session_ID: 1,
				}},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testRetrieveAllCases {

		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}

		defer func() {
			if err := logger.Sync(); err != nil {
				fmt.Println(err)
			}
		}()

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService}).Handle(w, r)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}

func TestDelete(t *testing.T) {
	testDeleteCases := []struct {
		name           string
		mockService    *test.MockService
		id             int64
		expectedStatus int
		prepare        func() *zap.Logger
	}{
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedResult: nil,
			},
			id:             15,
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			id:             15,
			expectedStatus: http.StatusUnprocessableEntity,
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

			w := httptest.NewRecorder()

			vars := map[string]string{
				"id": fmt.Sprintf("%d", tc.id),
			}

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%d", tc.id), nil)

			r = mux.SetURLVars(r, vars)

			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService}).HandleID(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}
