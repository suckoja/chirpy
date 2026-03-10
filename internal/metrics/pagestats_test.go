package metrics

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestNewPageStats(t *testing.T) {
	s := NewPageStats()
	if s == nil {
		t.Fatal("expected non-nil PageStats")
	}
	if s.GetHits() != 0 {
		t.Errorf("expected initial hits to be 0, got %d", s.GetHits())
	}
}

func TestIncrement(t *testing.T) {
	s := NewPageStats()
	s.Increment()
	s.Increment()
	s.Increment()
	if s.GetHits() != 3 {
		t.Errorf("expected 3 hits, got %d", s.GetHits())
	}
}

func TestReset(t *testing.T) {
	s := NewPageStats()
	s.Increment()
	s.Increment()
	s.Reset()
	if s.GetHits() != 0 {
		t.Errorf("expected 0 hits after reset, got %d", s.GetHits())
	}
}

func TestCountHitsMiddleware(t *testing.T) {
	s := NewPageStats()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := s.CountHits(next)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	if s.GetHits() != 5 {
		t.Errorf("expected 5 hits, got %d", s.GetHits())
	}
}

func TestCountHitsMiddlewareDelegates(t *testing.T) {
	s := NewPageStats()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := s.CountHits(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Errorf("expected next handler to be called, got status %d", w.Code)
	}
}

func TestIncrementConcurrent(t *testing.T) {
	s := NewPageStats()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.Increment()
		}()
	}
	wg.Wait()

	if s.GetHits() != goroutines {
		t.Errorf("expected %d hits, got %d", goroutines, s.GetHits())
	}
}
