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
	"github.com/darkjedidj/cinema-service/api/tickets"
	"github.com/darkjedidj/cinema-service/api/users"
)

type App struct {
	Router *mux.Router
}

// New creates router with handler
func (a *App) New(db *sql.DB, l *zap.Logger) {

	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/v1/tickets/{id}", tickets.Init(db, l).HandleID)
	myRouter.HandleFunc("/v1/tickets/{id}/download", users.Init(db, l).CheckTicket(tickets.Init(db, l).Download))
	myRouter.HandleFunc("/v1/tickets", users.Init(db, l).CheckPrivileges("tickets", tickets.Init(db, l).Handle))
	myRouter.HandleFunc("/v1/sessions/{id}/tickets", tickets.Init(db, l).Create)
	myRouter.HandleFunc("/v1/sessions/{id}", users.Init(db, l).CheckPrivileges("sessions", sessions.Init(db, l).HandleID))
	myRouter.HandleFunc("/v1/sessions", users.Init(db, l).CheckPrivileges("sessions", sessions.Init(db, l).Handle))
	myRouter.HandleFunc("/v1/halls/{id}/sessions", users.Init(db, l).CheckPrivileges("sessions", sessions.Init(db, l).Create))
	myRouter.HandleFunc("/v1/movies/{id}", users.Init(db, l).CheckPrivileges("movies", movies.Init(db, l).HandleID))
	myRouter.HandleFunc("/v1/movies", users.Init(db, l).CheckPrivileges("movies", movies.Init(db, l).Handle))
	myRouter.HandleFunc("/v1/halls/{id}", users.Init(db, l).CheckPrivileges("halls", halls.Init(db, l).HandleID))
	myRouter.HandleFunc("/v1/halls", users.Init(db, l).CheckPrivileges("halls", halls.Init(db, l).Handle))
	myRouter.HandleFunc("/v1/signin", users.Init(db, l).Signin)
	myRouter.HandleFunc("/v1/signup", users.Init(db, l).Signup)
	myRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8085/swagger/doc.json"), //The url pointing to API definition
	))
	a.Router = myRouter
}

// Run starts server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
