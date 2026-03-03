package app

import "github.com/suckoja/chirpy/internal/metrics"

type Server struct {
	stats *metrics.PageStats
}

func NewServer(stats *metrics.PageStats) *Server {
	return &Server{stats: stats}
}
