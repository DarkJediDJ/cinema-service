package tickets

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal"
	repo "github.com/darkjedidj/cinema-service/internal/repository/tickets"
	service "github.com/darkjedidj/cinema-service/internal/service/tickets"
	g "github.com/darkjedidj/cinema-service/package/generator"
)

type Handler struct {
	s   internal.Service // Allows use service features
	log *zap.Logger
	gen g.Client
}

func Init(db *sql.DB, l *zap.Logger) *Handler {

	service := service.Init(db, l)
	generator := g.Init(db, l)

	return &Handler{
		s:   service,
		log: l,
		gen: *generator,
	}
}

// HandleID handles all endpoints on this route
func (h *Handler) HandleID(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.Get(response, request) // GET BASE_URL/v1/tickets/{id}
	case http.MethodDelete:
		h.Delete(response, request) // DELETE BASE_URL/v1/tickets/{id}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handle handles all endpoints on this route
func (h *Handler) Handle(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		h.GetAll(response, request) // GET BASE_URL/v1/tickets
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Create get json and creates new ticket
// Create godoc
// @Security     ApiKeyAuth
// @Summary      Create ticket
// @Description  Creates ticket and returns created object
// @Tags      Tickets
// @Param        id    path  integer        true  "ticket ID"
// @Param        Body  body  repo.Resource  true  "The body to create a ticket"
// @Accept    json
// @Produce   json
// @Success      200  {object}  repo.Resource
// @Failure   400
// @Failure   422
// @Failure   500
// @Failure   401
// @Router       /sessions/{id}/tickets [post]
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var ticket repo.Resource

	response.Header().Set("Content-Type", "application/json")

	err = json.NewDecoder(request.Body).Decode(&ticket)
	if err != nil {
		h.log.Info("Failed to decode ticket json.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	ticket.Session_ID = int64(id)
	resource, err := h.s.Create(&ticket, ctx)
	if err != nil {

		if errors.Is(err, internal.ErrNoSeats) {
			response.WriteHeader(http.StatusBadRequest)

			_, err = response.Write([]byte(err.Error()))
			if err != nil {
				h.log.Info("Failed to write ticket response.",
					zap.Error(err),
				)

				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		if errors.Is(err, internal.ErrValidationFailed) {
			response.WriteHeader(http.StatusBadRequest)

			_, err = response.Write([]byte(err.Error()))
			if err != nil {
				h.log.Info("Failed to write ticket response.",
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
		h.log.Info("Failed to marshall ticket structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write ticket response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Delete get ID and deletes ticket with the same ID
// Delete godoc
// @Security     ApiKeyAuth
// @Summary      Delete ticket
// @Description  Deletes ticket
// @Param     id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /tickets/{id} [delete]
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
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

// Get ID and selects ticket with the same ID
// Get godoc
// @Security     ApiKeyAuth
// @Summary      Get ticket
// @Description  Gets ticket
// @Param        id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {object}  repo.Resource
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /tickets/{id} [get]
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
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
		h.log.Info("Failed to marshall ticket structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write ticket response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetAll selects all tickets
// GetAll godoc
// @Security     ApiKeyAuth
// @Summary      List ticket
// @Description  get tickets
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {array}  []repo.Resource
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router       /tickets [get]
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
		h.log.Info("Failed to marshall ticket structure.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = response.Write(body)
	if err != nil {
		h.log.Info("Failed to write ticket response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Download bought ticket
// Download godoc
// @Security  ApiKeyAuth
// @Summary   Download ticket
// @Param        id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success   200  {object}  g.Link
// @Failure      400
// @Failure      422
// @Failure      500
// @Failure      401
// @Router    /tickets/{id}/download [get]
func (h *Handler) Download(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.gen.GetTicket(ctx, int64(id))
	if err != nil {
		h.log.Info("Failed to get ticket fron bucket.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(url)
	if err != nil {
		h.log.Info("Failed to decode ticket json.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = response.Write(bf.Bytes())
	if err != nil {
		h.log.Info("Failed to write ticket response.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}
