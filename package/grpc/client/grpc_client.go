package generator

import (
	"context"
	"database/sql"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/darkjedidj/cinema-service/internal"
	t "github.com/darkjedidj/cinema-service/internal/repository/tickets"
	pb "github.com/darkjedidj/cinema-service/package/grpc/proto"
)

const port = "localhost:50051"

type Client struct {
	Repo *t.Repository
	Log  *zap.Logger
}

func Init(db *sql.DB, l *zap.Logger) *Client {

	return &Client{
		Repo: &t.Repository{DB: db, Log: l},
		Log:  l,
	}
}

// CreatePDF implements proto.NewTicketGeneratorServer
func (c *Client) CreatePDF(id int64, ctx context.Context) (int64, error) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return 0, internal.ErrInternalFailure
	}
	defer conn.Close()

	client := pb.NewTicketGeneratorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	entity, err := c.Repo.Retrieve(id, ctxTimeout)
	if err != nil {
		return 0, internal.ErrInternalFailure
	}

	res, ok := entity.(*t.Resource)
	if !ok {
		c.Log.Info("Failed to assert ticket object.",
			zap.Bool("ok", ok),
		)

		return 0, internal.ErrInternalFailure
	}

	result, err := client.GetTicket(ctx, &pb.TicketRequset{Time: res.Starts_at, Price: float32(res.Price), Seat: res.Seat, Id: res.ID, Title: res.Title})
	if err != nil {
		c.Log.Info("Failed to assert ticket object.",
			zap.Error(err),
		)
		return 0, internal.ErrInternalFailure
	}
	return result.ID, nil
}
