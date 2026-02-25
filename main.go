package main

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/metrics"
)

func main() {
	stats := &metrics.PageStats{}
	mux := routes(stats)

	// Listen and serve on port 8080 using the custom mux
	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

func routes(stats *metrics.PageStats) *http.ServeMux {
	mux := http.NewServeMux()

	// -- Static -- 
	// /app/*
	mux.Handle("/app/", stats.CountHits(mount("/app", http.Dir("."))))

	// /app/assets/*
	mux.Handle("/app/assets/", stats.CountHits(mount("/app/assets", http.Dir("./assets"))))

	// -- API --
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/metrics", metricsHandler(stats))
	mux.HandleFunc("POST /api/reset", resetMetricsHandler(stats))

	return mux
}

// Register the health check handler
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// Metrics Handler: Show page stat count
func metricsHandler(stats *metrics.PageStats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Hits: %d\n", stats.GetHits())
	}
}

// Reset Handler: Sets counter back to zero
func resetMetricsHandler(stats *metrics.PageStats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats.Reset()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hits reset to 0")
	}
}

func mount(prefix string, dir http.Dir) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(dir))
}
