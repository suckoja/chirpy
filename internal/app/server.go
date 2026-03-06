package app

import (
	"github.com/suckoja/chirpy/internal/database"
	"github.com/suckoja/chirpy/internal/metrics"
)

type Server struct {
	stats *metrics.PageStats
	db    *database.Queries
}

func NewServer(stats *metrics.PageStats, db *database.Queries) *Server {
	return &Server{stats: stats, db: db}
}
