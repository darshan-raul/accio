package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Git Service on :8082")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Git Service is healthy"))
	})
	log.Fatal(http.ListenAndServe(":8082", nil))
}
