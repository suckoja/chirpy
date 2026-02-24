package main

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/metrics"
)

func main() {
	mux := http.NewServeMux()

	// Ensure PageStats is initialized properly
	stats := &metrics.PageStats{}

	// Serve the main app at /app
	fileServer := http.FileServer(http.Dir("."))
	handleWithStats := stats.CountHits(http.StripPrefix("/app", fileServer))
	mux.Handle("/app/", handleWithStats)

	// Serve assets at /app/assets
	assetFileServer := http.FileServer(http.Dir("./assets"))
	assetHandle := http.StripPrefix("/app/assets", assetFileServer)
	mux.Handle("/app/assets/", assetHandle)

	// Register the health check handler
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// Metrics Handler: Show page stat count
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Hits: %d\n", stats.GetHits())
	})

	// Reset Handler: Sets counter back to zero
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		stats.Reset()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hits reset to 0")
	})

	// Listen and serve on port 8080 using the custom mux
	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}