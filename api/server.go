package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/darkjedidj/cinema-service/api/halls"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
	"go.uber.org/zap"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New(db *sql.DB, l *zap.Logger) {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/v1/halls/{id}", halls.Init(db, l).HandleID)
	myRouter.HandleFunc("/v1/halls", halls.Init(db, l).Handle)
	myRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8085/swagger/doc.json"), //The url pointing to API definition
	))
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
