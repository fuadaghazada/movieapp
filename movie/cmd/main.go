package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller"
	metadataGateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratingGateway "movieexample.com/movie/internal/gateway/rating/http"
	grpcHandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"net"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API Handler port")
	flag.Parse()

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

	metadataClient := metadataGateway.New(registry)
	ratingClient := ratingGateway.New(registry)
	ctrl := controller.New(ratingClient, metadataClient)

	////HTTP
	//h := httpHandler.New(ctrl)

	//http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	//if err := http.ListenAndServe(":8083", nil); err != nil {
	//	panic(err)
	//}

	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
}
