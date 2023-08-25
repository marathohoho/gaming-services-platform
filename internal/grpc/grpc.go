package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceRegistrar interface {
	RegisterGrpcServer(server *grpc.Server)
}

func CreateClientConnection(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil
	}

	return conn
}

func ListenForConnections(ctx context.Context, registrar ServiceRegistrar, addr, serviceName string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		// logrus.Fatal(err)
		return
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	registrar.RegisterGrpcServer(srv)

	// logrus.Infof("%s listening...", serviceName)

	go listenForStopped(ctx, srv, serviceName)

	if err = srv.Serve(lis); err != nil {
		// logrus.Fatal(err)
		return
	}
}

func listenForStopped(ctx context.Context, grpcServer *grpc.Server, serviceName string) {
	defer func() {
		// logrus.Infof("%s stopped", serviceName)
	}()

	for {
		select {
		case <-ctx.Done():
			grpcServer.Stop()
			return
		}
	}
}
