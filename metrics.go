package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		`<html>
  			<body>
    			<h1>Welcome, Chirpy Admin</h1>
    			<p>Chirpy has been visited %d times!</p>
  			</body>
		</html>`, cfg.fileHits.Load())))
}