package main

import "net/http"

func routes(s *Server) *http.ServeMux {
	mux := http.NewServeMux()

	// -- Static --
	mux.Handle("/app/", s.stats.CountHits(mount("/app", http.Dir("."))))
	mux.Handle("/app/assets/", s.stats.CountHits(mount("/app/assets", http.Dir("./assets"))))

	// -- API --
	mux.HandleFunc("GET /api/healthz", s.healthz)

	// -- Admin --
	mux.HandleFunc("GET /admin/metrics", s.metricsPage)
	mux.HandleFunc("POST /admin/reset", s.resetMetrics)

	return mux
}

func mount(prefix string, dir http.Dir) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(dir))
}