package api

import (
	context "context"
	"errors"
	"fmt"
	"net"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/services"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
)

type GrpcApi struct {
	config    *config.Config
	log       *logger.Logger
	svc       services.Service
	rate      *services.RateService
	rpcServer *grpc.Server
}

func NewGrpcApi(
	config *config.Config,
	svc services.Service,
	rate *services.RateService) *GrpcApi {

	middleware := services.NewLogMiddleware(config, svc)
	api := &GrpcApi{
		log:    logger.New(config.Log.Level),
		config: config,
		svc:    middleware,
		rate:   rate,
	}
	api.rpcServer = grpc.NewServer(grpc.UnaryInterceptor(api.rateLimiterInterceptor))

	return api
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

func (api *GrpcApi) rateLimiterInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	ip := md["client-ip"]
	if len(ip) > 0 {
		api.log.Info(fmt.Sprintf("cient ip %v", ip))
		valid, err := api.rate.Validate(
			ctx,
			ip[0],
			api.config.MaxRequests,
			(api.config.RateLimiter.WindowTimeInSeconds / 60))

		if err != nil {
			return nil, status.Errorf(codes.Unknown, "error unknown")
		}

		if !valid {
			return nil, status.Errorf(codes.ResourceExhausted, "Too many requests")
		}
	}

	return handler(ctx, req)
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
