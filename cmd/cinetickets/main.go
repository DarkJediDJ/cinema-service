package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	server "github.com/darkjedidj/cinema-service/api"
	"github.com/darkjedidj/cinema-service/api/halls"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = ":8085"

func main() {
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalln(err)
	}

	halls.Repo.DB = db

	a := server.App{}
	a.New()

	a.Run(port)
}
