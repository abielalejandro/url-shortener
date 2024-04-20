package api

import (
	"context"
	"net"

	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/services"
	"github.com/abielalejandro/tgs-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcApi struct {
	config    *config.Config
	log       *logger.Logger
	svc       services.Service
	rpcServer *grpc.Server
}

func NewRpcApi(config *config.Config, svc services.Service) *RpcApi {
	middleware := services.NewLogMiddleware(config, svc)
	return &RpcApi{
		rpcServer: grpc.NewServer(),
		log:       logger.New(config.Log.Level),
		config:    config,
		svc:       middleware,
	}
}

func (rpcApi *RpcApi) Run() {
	listener, err := net.Listen("tcp", rpcApi.config.RPC.Port)

	if err != nil {
		rpcApi.log.Fatal(err)
	}

	RegisterTgsServiceServer(rpcApi.rpcServer, rpcApi)
	reflection.Register(rpcApi.rpcServer)
	if err := rpcApi.rpcServer.Serve(listener); err != nil {
		rpcApi.log.Fatal(err)
	}
}

func (rpcApi *RpcApi) Next(context.Context, *NextRequest) (*NextResponse, error) {
	token, err := rpcApi.svc.GenerateToken()

	if err != nil {
		rpcApi.log.Fatal(err)
	}

	return &NextResponse{Token: token}, nil

}

func (rpcApi *RpcApi) Health(context.Context, *NextRequest) (*NextResponse, error) {
	return &NextResponse{Token: "OK"}, nil
}

func (rpcApi *RpcApi) mustEmbedUnimplementedTgsServiceServer() {
	rpcApi.log.Info("mustEmbedUnimplementedTgsServiceServer")
}
