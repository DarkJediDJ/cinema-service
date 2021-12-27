package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	config "github.com/darkjedidj/cinema-service/internal"
	db "github.com/darkjedidj/cinema-service/internal/queries"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Registration ...
func Registration(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user config.User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	db.CreateUser(user)
}
func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
