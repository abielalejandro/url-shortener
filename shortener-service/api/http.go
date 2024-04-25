package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	"github.com/abielalejandro/shortener-service/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CreateShortRequest struct {
	Url string `validate:"required,http_url"`
}

type HttpApi struct {
	config *config.Config
	log    *logger.Logger
	Router *mux.Router
	svc    services.Service
	rate   *services.RateService
}

type GeneralResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type ApiErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func sendResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, err string) {
	errPayload := &ApiErrorResponse{
		Code:    statusCode,
		Message: err,
		Errors:  []string{err},
	}

	payload := &GeneralResponse{Success: false, Data: errPayload, Timestamp: time.Now().UnixMilli()}
	sendResponse(w, statusCode, payload)
}

func NewHttpApi(
	config *config.Config,
	svc services.Service,
	rate *services.RateService) *HttpApi {

	middleware := services.NewLogMiddleware(config, svc)
	return &HttpApi{
		Router: mux.NewRouter().StrictSlash(true),
		log:    logger.New(config.Log.Level),
		config: config,
		svc:    middleware,
		rate:   rate,
	}
}

func (httpApi *HttpApi) handleRoutesV1() {
	subrouter := httpApi.Router.PathPrefix("/api/v1").Subrouter()
	subrouter.Use(httpApi.rateLimiterMiddleware)
	subrouter.HandleFunc("/short", httpApi.createShort).Methods("POST")
	subrouter.HandleFunc("/short/{id}", httpApi.searchUrlByShort).Methods("GET")
}

func (httpApi *HttpApi) health(w http.ResponseWriter, r *http.Request) {
	payload := &GeneralResponse{Success: true, Data: "UP", Timestamp: time.Now().UnixMilli()}
	sendResponse(w, http.StatusOK, payload)
}

func (httpApi *HttpApi) createShort(w http.ResponseWriter, r *http.Request) {
	var body CreateShortRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		httpApi.log.Error(err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		sendErrorResponse(w, http.StatusBadRequest, validationErrors.Error())
		return
	}

	token, err := httpApi.svc.GenerateShort(body.Url)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	payload := &GeneralResponse{Success: true, Data: token, Timestamp: time.Now().UnixMilli()}
	sendResponse(w, http.StatusOK, payload)

}

func (httpApi *HttpApi) searchUrlByShort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	token, err := httpApi.svc.SearchUrlByShort(id)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	payload := &GeneralResponse{Success: true, Data: token, Timestamp: time.Now().UnixMilli()}
	sendResponse(w, http.StatusOK, payload)

}

func (httpApi *HttpApi) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := utils.ReadUserIP(r)
		httpApi.log.Info(fmt.Sprintf("client ip %v", ip))
		valid, err := httpApi.rate.Validate(
			context.Background(),
			ip,
			httpApi.config.RateLimiter.MaxRequests,
			(httpApi.config.RateLimiter.WindowTimeInSeconds / 60))

		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if valid == false {
			sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (httpApi *HttpApi) Run() {
	httpApi.Router.HandleFunc("/health", httpApi.health).Methods("GET")
	httpApi.handleRoutesV1()
	httpApi.log.Fatal(http.ListenAndServe(httpApi.config.HTTP.Port, httpApi.Router))
}
