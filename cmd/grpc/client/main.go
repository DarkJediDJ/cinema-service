package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/darkjedidj/cinema-service/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewPDFGeneratorClient(conn)

	// Contact the server and print out its response.
	movie := "Matrix"
	if len(os.Args) > 1 {
		movie = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Generate(ctx, &proto.GenerateRequest{Movie: movie})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetFile())
}
