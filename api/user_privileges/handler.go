package user_privileges

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	_ "github.com/darkjedidj/cinema-service/docs"
	"github.com/darkjedidj/cinema-service/internal"
	repo "github.com/darkjedidj/cinema-service/internal/repository/user_privileges"
	service "github.com/darkjedidj/cinema-service/internal/service/user_privileges"
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
		h.Get(response, request) // GET BASE_URL/v1/user_privileges/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/user_privileges/{id}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handle handles all endpoints on this route
func (h *Handler) Handle(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.GetAll(response, request) // GET BASE_URL/v1/user_privileges
	case http.MethodPost:
		h.Create(response, request) // GET BASE_URL/v1/user_privileges
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Create get json and creates new User Privilege
// Create godoc
// @Security     ApiKeyAuth
// @Summary      Create User Privilege
// @Description  Creates User Privilege and returns created object
// @Tags         User Privileges
// @Param        Body  body  repo.Resource  true  "The body to create a User Privilege"
// @Accept       json
// @Produce      json
// @Success      200  {object}  repo.Resource
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /user_privileges [post]
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var user_privilege repo.Resource

	response.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&user_privilege)
	if err != nil {
		h.log.Info("Failed to decode user_privilege json.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	resource, err := h.s.Create(&user_privilege, ctx)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	body, err := json.Marshal(resource)
	if err != nil {
		h.log.Info("Failed to marshall user_privilege structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write user_privilege response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Delete get ID and deletes User Privilege with the same ID
// Delete godoc
// @Security     ApiKeyAuth
// @Summary      Delete User Privilege
// @Description  Deletes User Privilege
// @Param        id  path  integer  true  "User Privilege ID"
// @Tags         User Privileges
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /user_privileges/{id} [delete]
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse User Privilege id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.s.Delete(int64(id), ctx)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	response.WriteHeader(http.StatusOK)
}

// Get ID and selects User Privilege with the same ID
// Get godoc
// @Security     ApiKeyAuth
// @Summary      Get User Privilege
// @Description  Gets User Privilege
// @Param        id  path  integer  true  "User Privilege ID"
// @Tags         User Privileges
// @Accept       json
// @Produce      json
// @Success      200  {object}  repo.Resource
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /user_privileges/{id} [get]
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse User Privilege id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
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
		h.log.Info("Failed to marshall User Privilege structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write User Privilege response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetAll selects all User Privileges
// GetAll godoc
// @Security     ApiKeyAuth
// @Summary      List User Privileges
// @Description  get User Privileges
// @Tags         User Privileges
// @Accept       json
// @Produce      json
// @Success      200  {array}  []repo.Resource
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /user_privileges [get]
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
		h.log.Info("Failed to marshall User Privilege structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write User Privilege response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
