package main

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/metrics"
)

func main() {
	stats := &metrics.PageStats{}
	srv := NewServer(stats)

	mux := routes(srv)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}