package tckgenerator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/darkjedidj/cinema-service/internal/repository/tickets"
	cloud "github.com/darkjedidj/cinema-service/package/aws"
	g "github.com/darkjedidj/cinema-service/package/grpc/client"
)

type Link struct {
	URL string `json:"url"`
}

// Service is a struct to store DB and logger connection
type Client struct {
	gen *g.Client
	log *zap.Logger
}

// Init returns Service object
func Init(db *sql.DB, l *zap.Logger) *Client {

	return &Client{
		gen: &g.Client{Repo: &tickets.Repository{DB: db, Log: l}},
		log: l,
	}
}

// GetTicket creates PDF file, stores it in S3 and returns ID
func (c *Client) GetTicket(ctx context.Context, id int64) (*Link, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session := cloud.AwsService.GetSession()

	S3BucketName := os.Getenv("BUCKET_NAME")

	ticket, err := c.gen.CreatePDF(id, ctx)
	if err != nil {
		c.log.Info("Failed to assert ticket object.",
			zap.Error(err),
		)
		return nil, err
	}

	svc := s3.New(session)

	params := &s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(fmt.Sprintf("%v", ticket)),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(15 * time.Minute) // Set link expiration time
	if err != nil {
		c.log.Info("Failed to assert ticket object.",
			zap.Error(err),
		)
		return nil, err
	}

	fmt.Println(url)

	return &Link{URL: url}, nil
}
