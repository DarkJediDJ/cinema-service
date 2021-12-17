package authenticate

import (
	"log"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	AddHalls    bool
	AddSessions bool
	AddMovies   bool
}

//Registration ...
func Registration(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	db, err := sqlx.Connect("postgres", "user=admin dbname=db passwword=admin sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO user (login, password, addHalls, addSessions, addMovies) VALUES ($1, $2, $3,$4,$5)", "123", "1234", "False", "False", "False")
	response.Write([]byte("success"))
	u := User{}
	err = db.Get(&u, "SELECT * FROM person WHERE login=$1", "123")
	if err != nil {
		log.Fatalln(err)
	}
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
