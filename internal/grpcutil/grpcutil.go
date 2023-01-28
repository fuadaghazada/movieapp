package grpcutil

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"movieexample.com/pkg/discovery"
)

func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addr, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(addr[rand.Intn(len(addr))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
