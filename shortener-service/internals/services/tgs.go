package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	pb "github.com/abielalejandro/tgs-service/api"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TgsResponse struct {
	Success   bool   `json:"success"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type TgsService interface {
	Next(url string) (string, error)
}

type TgsServiceRest struct {
	config *config.Config
}
type TgsServiceGeneric struct {
	config *config.Config
}
type TgsServiceGrpc struct {
	config *config.Config
}

func (svc *TgsServiceRest) Next(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, svc.config.TgsService.Url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-URL-HASH", utils.ToBase62(url))
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var body TgsResponse
	err = json.Unmarshal(responseData, &body)
	if err != nil {
		return "", err
	}

	return body.Data, nil
}

func (svc *TgsServiceGeneric) Next(url string) (string, error) {
	return utils.ToBase62(url), nil
}

func (svc *TgsServiceGrpc) Next(url string) (string, error) {

	conn, err := grpc.Dial(svc.config.TgsService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := pb.NewTgsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Next(ctx, &pb.NextRequest{})
	if err != nil {
		return "", err
	}

	return r.Token, nil
}

func newTgsServiceRest(config *config.Config) *TgsServiceRest {
	return &TgsServiceRest{
		config: config,
	}
}

func newTgsServiceGeneric(config *config.Config) *TgsServiceGeneric {
	return &TgsServiceGeneric{config: config}
}

func newTgsServiceGrpc(config *config.Config) *TgsServiceGrpc {
	return &TgsServiceGrpc{config: config}
}

func NewTgsService(config *config.Config) TgsService {

	switch config.TgsService.Type {
	case "http":
		return newTgsServiceRest(config)
	case "grpc":
		return newTgsServiceGrpc(config)
	default:
		return newTgsServiceGeneric(config)
	}
}
