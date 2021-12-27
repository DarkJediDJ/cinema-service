package server

import (
	"log"
	"net/http"

	"github.com/darkjedidj/cinema-service/cmd/API/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New() {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/registration", handlers.Registration)
	myRouter.HandleFunc("/authentification", handlers.Authentification)
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
