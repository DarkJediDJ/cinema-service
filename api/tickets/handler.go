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
	ctx context.Context
	gen g.Client
}

func Init(db *sql.DB, l *zap.Logger, c context.Context) *Handler {

	service := service.Init(db, l)
	generator := g.Init(db, l)

	return &Handler{
		s:   service,
		log: l,
		ctx: c,
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
// @Summary      Create ticket
// @Description  Creates ticket and returns created object
// @Tags         Tickets
// @Param        id  path  integer  true  "ticket ID"
// @Param        Body  body  internal.Identifiable  true  "The body to create a ticket"
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /sessions/{id}/tickets [post]
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
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

	ticket.Session_ID = int64(id)
	resource, err := h.s.Create(&ticket, h.ctx)
	if err != nil {
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
// @Summary      Delete ticket
// @Description  Deletes ticket
// @Param        id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /tickets/{id} [delete]
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	err = h.s.Delete(int64(id), h.ctx)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	response.WriteHeader(http.StatusOK)
}

// Get ID and selects ticket with the same ID
// Get godoc
// @Summary      Get ticket
// @Description  Gets ticket
// @Param        id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /tickets/{id} [get]
func (h *Handler) Get(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	resource, err := h.s.Retrieve(int64(id), h.ctx)
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
// @Summary      List ticket
// @Description  get tickets
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {array}  []internal.Identifiable
// @Router       /tickets [get]
func (h *Handler) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	resource, err := h.s.RetrieveAll(h.ctx)
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
// @Summary      Download ticket
// @Param        id  path  integer  true  "ticket ID"
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Success      200  {object}  internal.Identifiable
// @Router       /tickets/{id}/dowload [get]
func (h *Handler) Download(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.log.Info("Failed to parse ticket id.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusBadGateway)
		return
	}

	url, err := h.gen.GetTicket(h.ctx, int64(id))
	if err != nil {
		h.log.Info("Failed to get ticket fron bucket.",
			zap.Error(err),
		)

		response.WriteHeader(http.StatusUnprocessableEntity)
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(url)

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
