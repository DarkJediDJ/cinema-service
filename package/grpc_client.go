package service

import (
	"context"
	"log"
	"time"

	"github.com/darkjedidj/cinema-service/internal"
	h "github.com/darkjedidj/cinema-service/internal/repository/halls"
	pb "github.com/darkjedidj/cinema-service/package/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = "localhost:50051"

type Client struct {
	repo *h.Repository
	log  *zap.Logger
}

// CreatePDF implements proto.NewTicketGeneratorServer
func (c *Client) CreatePDF(id int64) (int64, error) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTicketGeneratorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	entity, err := c.repo.Retrieve(id)
	if err != nil {
		return 0, internal.ErrInternalFailure
	}

	res, ok := entity.(*h.Resource)
	if !ok {
		c.log.Info("Failed to assert movie object.",
			zap.Bool("ok", ok),
		)

		return 0, internal.ErrInternalFailure
	}

	_, err = client.GetTicket(ctx, &pb.TicketRequset{res.Seats})
	if err != nil {
		log.Fatalf("could not generate ticket: %v", err)
	}
}
