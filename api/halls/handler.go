package halls

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/darkjedidj/cinema-service/internal"
	"go.uber.org/zap"

	"github.com/gorilla/mux"

	repo "github.com/darkjedidj/cinema-service/internal/repository/halls"
	service "github.com/darkjedidj/cinema-service/internal/service/halls"
)

type Handler struct {
	s   internal.Service // Allows use service features
	log *zap.Logger
}

func Init(db *sql.DB, l *zap.Logger) *Handler {

	service := service.Init(db, l)

	return &Handler{
		s:   service,
		log: l,
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

	response.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		h.log.Info("Failed to decode hall json.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	resource, err := h.s.Create(&hall)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		h.log.Info("Failed to marshall hall structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write hall response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Delete get ID and deletes Hall with the same ID
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse hall id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	err = h.s.Delete(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	response.WriteHeader(http.StatusOK)
}

// Get ID and selects Hall with the same ID
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse hall id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	resource, err := h.s.Retrieve(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if resource == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		h.log.Info("Failed to marshall hall structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write hall response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetAll selects all Halls
func (h *Handler) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	resource, err := h.s.RetrieveAll()
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if resource == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		h.log.Info("Failed to marshall hall structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write hall response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
