package main

import (
	"log"
	"net"

	server "github.com/darkjedidj/cinema-service/cmd/grpc/server/generator"
	"github.com/darkjedidj/cinema-service/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	srv := &server.GRPCServer{}
	proto.RegisterPDFGeneratorServer(s, srv)

	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
