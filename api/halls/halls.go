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

//CreateHall get json and creates new Hall
func CreateHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	Repo.Create(hall)
}

//DeleteHall get json and deletes Hall with the same ID
func DeleteHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	Repo.Delete(int64(hall.ID))
}

//SelectHall get json and selects Hall with the same ID
func SelectHall(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var hall halldb.Resource
	json.NewDecoder(request.Body).Decode(&hall)
	dbhall := Repo.Retrieve(int64(hall.ID))
	response.Write([]byte(fmt.Sprint(dbhall.ID)))
}
