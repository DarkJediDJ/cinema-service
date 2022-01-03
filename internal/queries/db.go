package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	config "github.com/darkjedidj/cinema-service/internal"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectToDb() *sql.DB {
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CreateUser(user config.User) {
	db := connectToDb()
	defer db.Close()
	insertUser := `insert into "User"("Login", "Password","AddHalls","AddMovies","AddSessions") values($1, $2,$3,$4,$5)`
	_, err := db.Exec(insertUser, user.Login, user.Password, user.AddHalls, user.AddMovies, user.AddSessions)
	if err != nil {
		panic(err)
	}
}

func SelectUser(user config.User) (dbUser config.User) {
	db := connectToDb()
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM "User" WHERE "Login" = $1`, user.Login)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&dbUser.UserID, &dbUser.Login, &dbUser.Password, &dbUser.AddHalls, &dbUser.AddMovies, &dbUser.AddSessions); err != nil {
			panic(err)
		}
	}
	return
}
