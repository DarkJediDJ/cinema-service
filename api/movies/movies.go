package movies

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	repo "github.com/darkjedidj/cinema-service/internal/repository/movies"
	service "github.com/darkjedidj/cinema-service/internal/service/movies"
)

type Handler struct {
	s *service.Service // Allows use service features
}

func Init(db *sql.DB) *Handler {

	service := service.Init(db)

	return &Handler{
		s: service,
	}
}

// HandleID handles all endpoints on this route
func (h *Handler) HandleID(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.Get(response, request) // GET BASE_URL/v1/movies/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/movies/{id}
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Handle handles all endpoints on this route
func (h *Handler) Handle(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.GetAll(response, request) // GET BASE_URL/v1/movies
	case http.MethodPost:
		h.Create(response, request) // POST BASE_URL/v1/movies
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Create get json and creates new Movie
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var movie repo.Resource

	err := json.NewDecoder(request.Body).Decode(&movie)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	dbmovie, err := h.s.Create(movie)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbmovie)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}

}

// Delete get ID and deletes Movie with the same ID
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	err = h.s.Delete(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	response.WriteHeader(http.StatusOK)
}

// Get ID and selects Movie with the same ID
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	dbhall, err := h.s.Retrieve(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbhall)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAll selects all Movie
func (h *Handler) GetAll(response http.ResponseWriter, request *http.Request) {

	dbhalls, err := h.s.RetrieveAll()
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbhalls)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}
