package main

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/metrics"
)

type Server struct {
	stats *metrics.PageStats
}

func NewServer(stats *metrics.PageStats) *Server {
	return &Server{stats: stats}
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func (s *Server) metricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Hits int
	}{
		Hits: s.stats.GetHits(),
	}

	if err := metricTpl.Execute(w, data); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (s *Server) resetMetrics(w http.ResponseWriter, r *http.Request) {
	s.stats.Reset()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hits reset to 0")
}