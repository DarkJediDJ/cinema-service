package main

import (
	server "github.com/darkjedidj/cinema-service/api"
	_ "github.com/lib/pq"
)

const port = ":8085"

func main() {

	a := server.App{}
	db := a.ConnectDB()
	defer db.Close()
	a.New(db)
	a.Run(port)
}
