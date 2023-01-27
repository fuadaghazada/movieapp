package main

import (
	"log"
	"movieexample.com/rating/internal/controller"
	httpHandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := controller.New(repo)
	h := httpHandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
