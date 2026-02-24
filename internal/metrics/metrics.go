package metrics

import (
	"net/http"
	"sync"
)

// HitCounter holds the in-memory state for user hits.
type PageStats struct {
	mu sync.Mutex // Protects the map from concurrent access
	hits int
}

func (s *PageStats) Increment() {
	s.mu.Lock()
	s.hits++
	s.mu.Unlock()
}

// GetHits returns the current count safely
func (s *PageStats) GetHits() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits
}

func (s *PageStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits = 0
}

// middleware function
func (s *PageStats) CountHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Increment()
		next.ServeHTTP(w, r)
	})
}