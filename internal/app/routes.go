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

	// -- /login
	mux.HandleFunc("POST /api/login", h.LogIn)

	// -- /users --
	mux.HandleFunc("POST /api/users", h.CreateUser)

	// -- /chirps --
	mux.HandleFunc("GET /api/chirps", h.ListChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", h.GetChirp)
	mux.HandleFunc("POST /api/chirps", h.CreateChirp)

	// -- Admin --
	mux.HandleFunc("GET /admin/metrics", h.MetricsPage)
	mux.HandleFunc("POST /admin/reset", h.ResetAll)

	return mux
}

func mount(prefix string, dir http.Dir) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(dir))
}
