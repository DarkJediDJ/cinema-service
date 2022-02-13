package test

import (
	"testing"

	hall "github.com/darkjedidj/cinema-service/internal/repository/halls"
)

func TestCreate(t *testing.T) {
	testCreateCases := []struct {
		name        string
		mockService MockService
	}{
		{name: "Success", mockService: MockService{ExpectedResult: hall.Repository{ID: 15, VIP: true, Seats: 15}, ExpectedError: nil, ExpectedArray: nil}},
	}
}

