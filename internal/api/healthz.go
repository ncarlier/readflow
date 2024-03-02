package api

import (
	"net/http"
	"sync/atomic"
)

var (
	healthy int32
)

// Shutdown set API as stopped
func Shutdown() {
	atomic.StoreInt32(&healthy, 0)
}

// Start set API as started
func Start() {
	atomic.StoreInt32(&healthy, 1)
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
