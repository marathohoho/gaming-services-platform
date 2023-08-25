package grpc

import (
	"gaming-services-platform/internal/grpc"
	"gaming-services-platform/proto"
)

func NewClient(address string) proto.BalanceServiceClient {
	connection := grpc.CreateClientConnection(address)
	return proto.NewBalanceServiceClient(connection)
}
