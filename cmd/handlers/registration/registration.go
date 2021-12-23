package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Registration ...
func Registration(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))

	insertUser := `insert into "User"("Login", "Password","AddHalls","AddMovies","AddSessions") values($1, $2,$3,$4,$5)`
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin dbname=admin password=admin sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	_, err = db.Exec(insertUser, user.Login, user.Password, user.AddHalls, user.AddMovies, user.AddSessions)
	if err != nil {
		panic(err)
	}
}
func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
