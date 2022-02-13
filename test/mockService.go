package test

import "github.com/darkjedidj/cinema-service/internal"

type MockService struct {
	ExpectedError  error
	ExpectedResult internal.Identifiable
}

// Create logic layer for repository method
func (s *MockService) Create(r internal.Identifiable) (internal.Identifiable, error) {
	return s.ExpectedResult, s.ExpectedError
}
