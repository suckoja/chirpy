package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) ResetAll(w http.ResponseWriter, r *http.Request) {
	if err := h.db.DeleteAllUsers(r.Context()); err != nil {
		http.Error(w, "failed to delete all users", http.StatusInternalServerError)
		return
	}

	h.stats.Reset()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hits reset to 0 and all users deleted")
}
