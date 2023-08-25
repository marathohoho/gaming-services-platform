package grpc

import (
	"context"
	"errors"
	"gaming-services-platform/internal/grpc"
	"gaming-services-platform/internal/repositories"
	"gaming-services-platform/proto"

	"github.com/go-redis/redis/v8"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "wallet grpc service"

type GrpcService struct {
	rdb           *redis.Client
	repository    *repositories.WalletRepository
	serverAddress string
	proto.UnimplementedBalanceServiceServer
}

func (gs *GrpcService) ListenForConnection(ctx context.Context) {
	grpc.ListenForConnections(ctx, gs, gs.serverAddress, serviceName)
}

func (gs *GrpcService) RegisterGrpcServer(server *grpcLib.Server) {
	proto.RegisterBalanceServiceServer(server, gs)
}

func (gs *GrpcService) Get(ctx context.Context, req *proto.RequestBalance) (*proto.ResponseBalance, error) {
	userWallet, err := gs.repository.Get(req.UserId)
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("unable to find a wallet by userId provided")
		}
		return nil, err
	}

	return &proto.ResponseBalance{Amount: userWallet.Amount}, nil
}
func NewGrpcService(redis *redis.Client, walletRepository *repositories.WalletRepository, address string) *GrpcService {
	return &GrpcService{
		rdb:           redis,
		repository:    walletRepository,
		serverAddress: address,
	}
}
