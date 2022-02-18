package test

import (
	"net/http"
	"testing"

	"github.com/darkjedidj/cinema-service/api"
	hall "github.com/darkjedidj/cinema-service/internal/repository/halls"
)
var a = server.App{}

func TestCreate(t *testing.T) {
	testCreateCases := []struct {
		name        string
		mockService MockService
	}{
		{name: "Success", mockService: MockService{ExpectedResult: &hall.Resource{ID: 15, VIP: true, Seats: 15}, ExpectedError: nil, ExpectedArray: nil}},
	}
	for _, tc := range testCreateCases {
		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", "/generate/result", nil)
			if err != nil {
				t.Fatal(err)
			}

			

		})
	}
}
