package handlers

import (
	"github.com/suckoja/chirpy/internal/database"
	"github.com/suckoja/chirpy/internal/metrics"
)

type Handlers struct {
	stats *metrics.PageStats
	db    *database.Queries
	jwtSecret string
}

func New(stats *metrics.PageStats, db *database.Queries, jwtSecret string) *Handlers {
	return &Handlers{stats: stats, db: db, jwtSecret: jwtSecret}
}
