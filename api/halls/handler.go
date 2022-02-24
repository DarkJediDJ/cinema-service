package halls

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/darkjedidj/cinema-service/internal"

	"github.com/gorilla/mux"

	repo "github.com/darkjedidj/cinema-service/internal/repository/halls"
	service "github.com/darkjedidj/cinema-service/internal/service/halls"
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
		h.Get(response, request) // GET BASE_URL/v1/halls/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/halls/{id}
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Handle handles all endpoints on this route
func (h *Handler) Handle(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.GetAll(response, request) // GET BASE_URL/v1/halls
	case http.MethodPost:
		h.Create(response, request) // POST BASE_URL/v1/halls
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

// Create get json and creates new Hall
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {
	var hall repo.Resource

	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		// TODO logging ZAP

		response.WriteHeader(http.StatusBadGateway) // TODO StatusBadGateway ???... 400
		return
	}

	resource, err := h.s.Create(&hall)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		log.Fatal(err)
	}

	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(body)
	if err != nil {
		log.Fatal(err)
	}

}

// Delete get ID and deletes Hall with the same ID
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

// Get ID and selects Hall with the same ID
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	resource, err := h.s.Retrieve(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if resource == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(resource)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAll selects all Halls
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
