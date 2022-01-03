package server

import (
	"context"

	"github.com/darkjedidj/cinema-service/pkg/proto"
)

type GRPCServer struct {
	proto.UnimplementedPDFGeneratorServer
}

func (s *GRPCServer) Generate(ctx context.Context, req *proto.GenerateRequest) (*proto.GenerateResponse, error) {
	return &proto.GenerateResponse{File: req.Movie}, nil
}
