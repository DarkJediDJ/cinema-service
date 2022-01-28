package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/darkjedidj/cinema-service/api/halls"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New(db *sql.DB) {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("v1/halls/{id}", halls.Init(db).HandleID)
	myRouter.HandleFunc("v1/halls", halls.Init(db).Handle)
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
