package handlers

import (
	"fmt"
	"net/http"

	"github.com/suckoja/chirpy/internal/templates"
)

func (h *Handlers) MetricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Hits int
	}{
		Hits: h.stats.GetHits(),
	}

	if err := templates.MetricTpl.Execute(w, data); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (h *Handlers) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	h.stats.Reset()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hits reset to 0")
}