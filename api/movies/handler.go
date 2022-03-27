package movies

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	repo "github.com/darkjedidj/cinema-service/internal/repository/movies"
	service "github.com/darkjedidj/cinema-service/internal/service/movies"
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
		h.Get(response, request) // GET BASE_URL/v1/movies/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/movies/{id}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
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
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Create get json and creates new Movie
// Create godoc
// @Summary      Create movie
// @Description  Creates movie and returns created object
// @Tags         Movies
// @Param         Body  body  internal.Identifiable  true  "The body to create a movie"
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /movies [post]
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var movie repo.Resource

	response.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&movie)
	if err != nil {
		h.log.Info("Failed to decode movie json.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	resource, err := h.s.Create(&movie, ctx)
	if err != nil {
		if errors.Is(err, internal.ErrValidationFailed) {
			response.WriteHeader(http.StatusBadRequest)

			_, err = response.Write([]byte(err.Error()))
			if err != nil {
				h.log.Info("Failed to write movies response.",
					zap.Error(err),
				)

				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		response.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		h.log.Info("Failed to marshall movie structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write movie response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Delete get ID and deletes movie with the same ID
// Delete godoc
// @Summary      Delete movie
// @Description  Deletes movie
// @Param        id  path  integer  true  "Movie ID"
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /movies/{id} [delete]
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse movie id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	err = h.s.Delete(int64(id), ctx)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	response.WriteHeader(http.StatusOK)
}

// Get ID and selects movie with the same ID
// Get godoc
// @Summary      Get movie
// @Description  Gets movie
// @Param        id  path  integer  true  "Movie ID"
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /movies/{id} [get]
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse movie id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	resource, err := h.s.Retrieve(int64(id), ctx)
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
		h.log.Info("Failed to marshall movie structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write movie response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetAll selects all movies
// GetAll godoc
// @Summary      List movie
// @Description  get movies
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Success      200  {array}  []internal.Identifiable
// @Router       /movies [get]
func (h *Handler) GetAll(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response.Header().Set("Content-Type", "application/json")

	resource, err := h.s.RetrieveAll(ctx)
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
		h.log.Info("Failed to marshall movie structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write movie response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
