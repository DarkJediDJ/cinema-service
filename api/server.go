package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/api/halls"
	"github.com/darkjedidj/cinema-service/api/movies"
	"github.com/darkjedidj/cinema-service/api/sessions"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New(db *sql.DB, l *zap.Logger) {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/v1/sessions/{id}", sessions.Init(db, l).HandleID)
	myRouter.HandleFunc("/v1/sessions", sessions.Init(db, l).Handle)
	myRouter.HandleFunc("/v1/halls/{id}/sessions", sessions.Init(db, l).Create)
	myRouter.HandleFunc("/v1/movies/{id}", movies.Init(db, l).HandleID)
	myRouter.HandleFunc("/v1/movies", movies.Init(db, l).Handle)
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
