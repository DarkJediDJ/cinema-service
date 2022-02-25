package halls

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darkjedidj/cinema-service/internal"
	hall "github.com/darkjedidj/cinema-service/internal/repository/halls"
	"github.com/darkjedidj/cinema-service/test"
)

func TestCreate(t *testing.T) {
	testCreateCases := []struct {
		name           string
		mockService    *test.MockService
		body           string
		expectedStatus int
	}{
		{
			name: "failure: empty body",
			mockService: &test.MockService{
				ExpectedResult: &hall.Resource{ID: 15, VIP: true, Seats: 15},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedResult: &hall.Resource{ID: 15, VIP: true, Seats: 15},
			},
			body:           `{"VIP": true, "seats": 10}`,
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			body:           `{"VIP": true, "seats": 123}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCreateCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService}).Handle(w, r)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}

func TestRetrieve(t *testing.T) {
	testRetrieveCases := []struct {
		name           string
		mockService    *test.MockService
		body           string
		expectedStatus int
	}{
		{
			name: "failure: empty body",
			mockService: &test.MockService{
				ExpectedResult: &hall.Resource{ID: 15, VIP: true, Seats: 15},
			},
			expectedStatus: http.StatusBadGateway,
		},
		{
			name: "success",
			mockService: &test.MockService{
				ExpectedResult: &hall.Resource{ID: 15, VIP: true, Seats: 15},
			},
			body:           `{"VIP": true, "seats": 10}`,
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure: DB error",
			mockService: &test.MockService{
				ExpectedError: internal.ErrInternalFailure,
			},
			body:           `{"VIP": true, "seats": 123}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testRetrieveCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")

			(&Handler{s: tc.mockService}).Handle(w, r)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}
