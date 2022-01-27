package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/darkjedidj/cinema-service/api/halls"
	hall "github.com/darkjedidj/cinema-service/internal/repository/halls"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New(db *sql.DB) {

	myRouter := mux.NewRouter().StrictSlash(false)
	handler := halls.Handler{Repo: &hall.Repository{DB: db}}
	myRouter.HandleFunc("/v1/halls/{id}", handler.Handle)
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
