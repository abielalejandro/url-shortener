package api

import (
	context "context"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type TgsServiceGrpc struct {
	config *config.Config
}

func (svc *TgsServiceGrpc) Next(url string) (string, error) {

	conn, err := grpc.Dial(svc.config.TgsService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := NewTgsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "x-url-hash", url)

	r, err := c.Next(ctx, &NextRequest{})
	if err != nil {
		return "", err
	}

	return r.Token, nil
}

func NewTgsServiceGrpc(config *config.Config) *TgsServiceGrpc {
	return &TgsServiceGrpc{config: config}
}
