package handlers

import "github.com/suckoja/chirpy/internal/metrics"

type Handlers struct {
	stats *metrics.PageStats
}

func New(stats *metrics.PageStats) *Handlers {
	return &Handlers{stats: stats}
}
