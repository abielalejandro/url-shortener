package services

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/pkg/utils"
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
type TgsServiceGeneric struct{}

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

func newTgsServiceRest(config *config.Config) *TgsServiceRest {
	return &TgsServiceRest{
		config: config,
	}
}

func newTgsServiceGeneric(config *config.Config) *TgsServiceGeneric {
	return &TgsServiceGeneric{}
}

func NewTgsService(config *config.Config) TgsService {

	switch config.TgsService.Type {
	case "http":
		return newTgsServiceRest(config)
	default:
		return newTgsServiceGeneric(config)
	}
}
