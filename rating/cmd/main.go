package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"log"
	"movieexample.com/gen"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller"
	grpcHandler "movieexample.com/rating/internal/handler/grpc"
	"movieexample.com/rating/internal/repository/mysql"
	"net"
	"os"
	"time"
)

const serviceName = "rating"

func main() {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	port := cfg.API.Port

	log.Printf("Starting the %v service\n", serviceName)

	// registry
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("Failed to report healthy state: %v\n", err.Error())
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	//ingester, err := kafka.NewIngester("localhost", "", "")
	ctrl := controller.New(repo, nil)

	//// HTTP
	//h := httpHandler.New(ctrl)
	//
	//http.Handle("/rating", http.HandlerFunc(h.Handle))
	//if err := http.ListenAndServe(":8082", nil); err != nil {
	//	panic(err)
	//}

	h := grpcHandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}
