package main

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/metrics"
	"github.com/suckoja/chirpy/internal/app"
)

func main() {
	stats := &metrics.PageStats{}
	srv := app.NewServer(stats)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", srv.Routes()); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
