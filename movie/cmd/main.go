package main

import (
	"log"
	"movieexample.com/movie/internal/controller"
	metadataGateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratingGateway "movieexample.com/movie/internal/gateway/rating/http"
	httpHandler "movieexample.com/movie/internal/handler/http"
	"net/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataClient := metadataGateway.New("localhost:8081")
	ratingClient := ratingGateway.New("localhost:8082")
	ctrl := controller.New(ratingClient, metadataClient)
	h := httpHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
