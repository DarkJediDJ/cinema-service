package test

import (
	"context"

	"github.com/darkjedidj/cinema-service/internal"
)

type MockService struct {
	ExpectedError  error
	ExpectedResult internal.Identifiable
	ExpectedArray  []internal.Identifiable
}

func (s *MockService) Create(_ internal.Identifiable, _ context.Context) (internal.Identifiable, error) {
	return s.ExpectedResult, s.ExpectedError
}

func (s *MockService) Retrieve(_ int64, _ context.Context) (internal.Identifiable, error) {
	return s.ExpectedResult, s.ExpectedError
}

func (s *MockService) RetrieveAll(_ context.Context) ([]internal.Identifiable, error) {

	return s.ExpectedArray, s.ExpectedError
}

func (s *MockService) Delete(_ int64, _ context.Context) error {
	return s.ExpectedError
}
