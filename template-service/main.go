package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Template Service on :8081")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Template Service is healthy"))
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
