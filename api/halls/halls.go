package halls

import (
	"encoding/json"
	"fmt"
	"net/http"

	halldb "github.com/darkjedidj/cinema-service/internal/repository/halls"
	"github.com/gorilla/mux"
)

var Repo halldb.Repository

func Handle(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	route := vars["route"]
	switch route {
	case "add":
		CreateHall(response, request)
	case "get":
		SelectHall(response, request)
	case "delete":
		DeleteHall(response, request)
	default:
		response.WriteHeader(http.StatusBadGateway)
	}
}

func CreateHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	Repo.Create(hall)
}

func DeleteHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	Repo.Delete(hall)
}

func SelectHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	dbhall := Repo.Retrieve(int64(hall.ID))
	response.Write([]byte(fmt.Sprint(dbhall.ID)))
}
