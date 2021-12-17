package main

import (
	"log"
	"net/http"

	auth "github.com/darkjedidj/cinema-service/handlers/authenticate"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New() {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/registration", auth.Registration)
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func main() {
	a := App{}
	a.New()

	a.Run(":8085")
}
