package halls

import (
	"encoding/json"
	"fmt"
	"net/http"

	halldb "github.com/darkjedidj/cinema-service/internal/repository/halls"
	"github.com/gorilla/mux"
)

var Repo halldb.Repository

//Handle handles all endpoints on this route
func Handle(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	route := vars["route"]
	switch route {
	case "post":
		Post(response, request)
	case "get":
		Get(response, request)
	case "delete":
		Delete(response, request)
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

//Create get json and creates new Hall
func Post(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}
	err = Repo.Create(hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}
}

//Delete get json and deletes Hall with the same ID
func Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}
	err = Repo.Delete(int64(hall.ID))
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}
}

//Select get json and selects Hall with the same ID
func Get(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}
	dbhall, err := Repo.Retrieve(int64(hall.ID))
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}
	response.Write([]byte(fmt.Sprint(dbhall.ID)))
}
