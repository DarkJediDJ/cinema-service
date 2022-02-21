package sessions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/darkjedidj/cinema-service/internal"

	"github.com/gorilla/mux"

	repo "github.com/darkjedidj/cinema-service/internal/repository/sessions"
	service "github.com/darkjedidj/cinema-service/internal/service/sessions"
)

type Handler struct {
	s internal.Service // Allows use service features
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
		h.Get(response, request) // GET BASE_URL/v1/sessions/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/sessions/{id}
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Handle handles all endpoints on this route
func (h *Handler) Handle(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.GetAll(response, request) // GET BASE_URL/v1/sessions
	case http.MethodPost:
		h.Create(response, request) // POST BASE_URL/v1/sessions
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Create get json and creates new session
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var session repo.Resource

	err := json.NewDecoder(request.Body).Decode(&session)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	dbsession, err := h.s.Create(&session)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	res, err := json.Marshal(dbsession)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}

}

// Delete get ID and deletes session with the same ID
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

// Get ID and selects session with the same ID
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	dbsession, err := h.s.Retrieve(int64(id))
	if err != nil {
		fmt.Print(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbsession)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAll selects all session
func (h *Handler) GetAll(response http.ResponseWriter, request *http.Request) {

	dbsessions, err := h.s.RetrieveAll()
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbsessions)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}
