package halls

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	halldb "github.com/darkjedidj/cinema-service/internal/repository/halls"
	"github.com/gorilla/mux"
)

var Repo halldb.Repository

//Handle handles all endpoints on this route
func Handle(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			Create(response, request, db) //[POST] BASE_URL/v1/halls/create
		case http.MethodGet:
			Get(response, request, db) //[GET] BASE_URL/v1/halls/get
		case http.MethodDelete:
			Delete(response, request, db) //[DELETE] BASE_URL/v1/halls/delete
		default:
			response.WriteHeader(http.StatusBadGateway)
		}
	}
}

//Create get json and creates new Hall
func Create(response http.ResponseWriter, request *http.Request, db *sql.DB) {

	response.Header().Set("Content-Type", "application/json")

	var hall halldb.Resource

	err := json.NewDecoder(request.Body).Decode(&hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	Repo.DB = db

	err = Repo.Create(hall)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}
}

//Delete get json and deletes Hall with the same ID
func Delete(response http.ResponseWriter, request *http.Request, db *sql.DB) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	Repo.DB = db

	err = Repo.Delete(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}
}

//Get get json and selects Hall with the same ID
func Get(response http.ResponseWriter, request *http.Request, db *sql.DB) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
	}

	Repo.DB = db

	dbhall, err := Repo.Retrieve(int64(id))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(dbhall)
	if err != nil {
		log.Fatal(err)
	}
	response.Header().Set("Content-Type", "application/json")

	_, err = response.Write(res)
	if err != nil {
		log.Fatal(err)
	}
}
