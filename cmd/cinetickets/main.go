package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	server "github.com/darkjedidj/cinema-service/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = ":8085"

func main() {

	a := server.App{}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	a.New(db)
	a.Run(port)
}
