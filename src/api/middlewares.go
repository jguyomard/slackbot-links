package api

import "net/http"

// APIMiddleware adds some headers (CORS, Content-Type...)
type APIMiddleware struct {
	handler http.Handler
}

// NewAPIMiddleware create a new API middleware
func NewAPIMiddleware(handler http.Handler) *APIMiddleware {
	return &APIMiddleware{handler: handler}
}

func (m *APIMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	m.handler.ServeHTTP(w, r)
}
