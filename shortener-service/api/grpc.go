package api

import (
	context "context"
	"errors"
	"net"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
)

type GrpcApi struct {
	config    *config.Config
	log       *logger.Logger
	svc       services.Service
	rpcServer *grpc.Server
}

func NewGrpcApi(config *config.Config, svc services.Service) *GrpcApi {
	middleware := services.NewLogMiddleware(config, svc)
	return &GrpcApi{
		rpcServer: grpc.NewServer(),
		log:       logger.New(config.Log.Level),
		config:    config,
		svc:       middleware,
	}
}

func (api *GrpcApi) Run() {
	listener, err := net.Listen("tcp", api.config.GRPC.Port)

	if err != nil {
		api.log.Fatal(err)
	}

	RegisterShortenerServiceServer(api.rpcServer, api)
	reflection.Register(api.rpcServer)
	if err := api.rpcServer.Serve(listener); err != nil {
		api.log.Fatal(err)
	}
}

func (api *GrpcApi) Health(ctx context.Context, req *HealthRequest) (*HealthResponse, error) {
	return &HealthResponse{Msg: "UP"}, nil
}

func (api *GrpcApi) mustEmbedUnimplementedShortenerServiceServer() {
	api.log.Info("mustEmbedUnimplementedShortenerServiceServer")
}

func (api *GrpcApi) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	token, err := api.svc.SearchUrlByShort(req.Url)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return &SearchResponse{Url: token}, nil
}

func (api *GrpcApi) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	token, err := api.svc.GenerateShort(req.Url)
	if err != nil {
		if errors.Is(err, &storage.NotFoundError{}) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return &CreateResponse{Url: token}, nil
}
