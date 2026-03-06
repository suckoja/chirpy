package app

import (
	"net/http"

	"github.com/suckoja/chirpy/internal/handlers"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	h := handlers.New(s.stats, s.db)

	// -- Static --
	mux.Handle("/app/", s.stats.CountHits(mount("/app", http.Dir("."))))
	mux.Handle("/app/assets/", s.stats.CountHits(mount("/app/assets", http.Dir("./assets"))))

	// -- API --
	mux.HandleFunc("GET /api/healthz", h.Healthz)
	mux.HandleFunc("POST /api/validate_chirp", h.ValidateChirp)

	// example future route
	// mux.HandleFunc("POST /api/users", h.CreateUser)

	// -- Admin --
	mux.HandleFunc("GET /admin/metrics", h.MetricsPage)
	mux.HandleFunc("POST /admin/reset", h.ResetMetrics)

	return mux
}

func mount(prefix string, dir http.Dir) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(dir))
}
