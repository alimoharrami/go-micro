package main

import (
	"go-blog/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/post", handlers.ViewPostHandler)
	http.HandleFunc("/new", handlers.NewPostHandler)
	http.HandleFunc("/create", handlers.CreatePostHandler)
	http.ListenAndServe(":8080", nil)
}
