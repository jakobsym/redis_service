package main

import (
	"fmt"
	"net/http"
)

// can send curl commands to server to test if working as intended
func main() {
	server := &http.Server{
		Addr:    ":3000",
		Handler: http.HandlerFunc(basicHandler),
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("error listening to server")
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}
