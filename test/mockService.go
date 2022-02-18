package test

import (
	"github.com/darkjedidj/cinema-service/internal"
)

type MockService struct {
	ExpectedError  error
	ExpectedResult internal.Identifiable
	ExpectedArray  []internal.Identifiable
}

func (s *MockService) Create(r internal.Identifiable) (internal.Identifiable, error) {
	return s.ExpectedResult, s.ExpectedError
}

func (s *MockService) Retrieve(id int64) (internal.Identifiable, error) {
	return s.ExpectedResult, s.ExpectedError
}

func (s *MockService) RetrieveAll() ([]internal.Identifiable, error) {

	return s.ExpectedArray, s.ExpectedError
}

func (s *MockService) Delete(r internal.Identifiable) error {
	return s.ExpectedError
}
