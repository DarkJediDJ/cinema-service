package main

import (
	"github.com/darkjedidj/cinema-service/cmd/API/server"
)

func main() {
	a := server.App{}
	a.New()

	a.Run(":8085")
}
