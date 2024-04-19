package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/services"
	"github.com/abielalejandro/tgs-service/pkg/logger"
	"github.com/gorilla/mux"
)

type HttpApi struct {
	config *config.Config
	log    *logger.Logger
	Router *mux.Router
	svc    services.Service
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

func NewHttpApi(config *config.Config, svc services.Service) *HttpApi {
	middleware := services.NewLogMiddleware(config, svc)
	return &HttpApi{
		Router: mux.NewRouter().StrictSlash(true),
		log:    logger.New(config.Log.Level),
		config: config,
		svc:    middleware,
	}
}

func (httpApi *HttpApi) handleRoutesV1() {
	subrouter := httpApi.Router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/next", httpApi.getNextToken).Methods("GET")
}

func (httpApi *HttpApi) getNextToken(w http.ResponseWriter, r *http.Request) {
	token, err := httpApi.svc.GenerateToken()

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := &GeneralResponse{Success: true, Data: token, Timestamp: time.Now().UnixMilli()}
	sendResponse(w, http.StatusOK, payload)

}

func (httpApi *HttpApi) Run() {
	httpApi.handleRoutesV1()
	log.Fatal(http.ListenAndServe(httpApi.config.Port, httpApi.Router))
}
