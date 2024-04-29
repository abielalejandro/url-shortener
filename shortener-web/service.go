package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	api "github.com/abielalejandro/shortener-web/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShortenerService interface {
	Create(longUrl string, sourceIp string) (string, error)
	Search(short string, sourceIp string) (string, error)
	Health(sourceIp string) (string, error)
}

type HttpShortenerService struct {
	config *Config
}

type GrpcShortenerService struct {
	config *Config
}

type ShortenerRequest struct {
	Url string `json:"url"`
}

type ShortenerResponseErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ShortenerResponseError struct {
	Success   bool                       `json:"success"`
	Data      ShortenerResponseErrorData `json:"data"`
	Timestamp int64                      `json:"timestamp"`
}

type ShortenerResponse struct {
	Success   bool   `json:"success"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func NewShortenerService(
	config *Config,
) ShortenerService {

	switch config.RestShortenerService.Type {
	case "http":
		return &HttpShortenerService{
			config: config,
		}
	case "grpc":
		return &GrpcShortenerService{
			config: config,
		}
	default:
		return &HttpShortenerService{
			config: config,
		}
	}
}

// HttpShortenerService
func (svc *HttpShortenerService) Create(longUrl string, sourceIp string) (string, error) {
	client := &http.Client{}
	data := &ShortenerRequest{
		Url: longUrl,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v/api/%v/short", svc.config.RestShortenerService.Url, svc.config.RestShortenerService.Version)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Real-Ip", sourceIp)
	req.Header.Set("X-Forwarded-For", sourceIp)
	req.Header.Set("X-URL-HASH", ToBase62(sourceIp))
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		var body ShortenerResponseError
		err = json.Unmarshal(responseData, &body)
		if err != nil {
			return "", err
		}

		return "", errors.New(body.Data.Message)
	}

	var body ShortenerResponse
	err = json.Unmarshal(responseData, &body)
	if err != nil {
		return "", err
	}

	return body.Data, nil
}

func (svc *HttpShortenerService) Search(short string, sourceIp string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%v/api/%v/short/%v", svc.config.RestShortenerService.Url, svc.config.RestShortenerService.Version, short)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	req.Header.Set("X-Real-Ip", sourceIp)
	req.Header.Set("X-Forwarded-For", sourceIp)
	req.Header.Set("X-URL-HASH", ToBase62(sourceIp))
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		var body ShortenerResponseError
		err = json.Unmarshal(responseData, &body)
		if err != nil {
			return "", err
		}
		return "", errors.New(body.Data.Message)
	}

	var body ShortenerResponse
	err = json.Unmarshal(responseData, &body)
	if err != nil {
		return "", err
	}

	return body.Data, nil
}

func (svc *HttpShortenerService) Health(sourceIp string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%v/health", svc.config.RestShortenerService.Url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	req.Header.Set("X-Real-Ip", sourceIp)
	req.Header.Set("X-Forwarded-For", sourceIp)
	req.Header.Set("X-URL-HASH", ToBase62(sourceIp))
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		var body ShortenerResponseError
		err = json.Unmarshal(responseData, &body)
		if err != nil {
			return "", err
		}
		return "", errors.New(body.Data.Message)
	}

	var body ShortenerResponse
	err = json.Unmarshal(responseData, &body)
	if err != nil {
		return "", err
	}

	return body.Data, nil
}

// GrpcShortenerService
func (svc *GrpcShortenerService) Create(longUrl string, sourceIp string) (string, error) {

	conn, err := grpc.Dial(svc.config.RestShortenerService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := api.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Create(ctx, &api.CreateRequest{
		Url: longUrl,
	})

	if err != nil {
		return "", err
	}

	return r.Url, nil
}

func (svc *GrpcShortenerService) Search(short string, sourceIp string) (string, error) {

	conn, err := grpc.Dial(svc.config.RestShortenerService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := api.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Search(ctx, &api.SearchRequest{
		Url: short,
	})

	if err != nil {
		return "", err
	}

	return r.Url, nil

}

func (svc *GrpcShortenerService) Health(sourceIp string) (string, error) {

	conn, err := grpc.Dial(svc.config.RestShortenerService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := api.NewShortenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Health(ctx, &api.HealthRequest{})

	if err != nil {
		return "", err
	}

	return r.Msg, nil

}
