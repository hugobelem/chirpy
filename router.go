package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/hugobelem/chirpy/internal/database"
)

func setupRouter() *http.ServeMux {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error openning database: %s", err)
	}
	dbQueries := database.New(db)

	config := apiConfig{
		fileHits: atomic.Int32{},
		db:       dbQueries,
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
	}

	fileHandler := http.StripPrefix(filepath, http.FileServer(http.Dir(dir)))
	mux := http.NewServeMux()
	mux.Handle(filepath, config.middlewareMetrics(fileHandler))
	mux.HandleFunc("GET  /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", config.handlerCreateUser)
	mux.HandleFunc("PUT  /api/users", config.handlerUpdateUser)
	mux.HandleFunc("POST /api/chirps", config.handlerCreateChirps)
	mux.HandleFunc("GET  /api/chirps", config.handlerRetrieveChirps)
	mux.HandleFunc("GET  /api/chirps/{chirpID}", config.handlerGetSingleChirp)
	mux.HandleFunc("POST /api/login", config.handlerLogin)
	mux.HandleFunc("POST /api/refresh", config.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", config.handlerRevokeToken)

	mux.HandleFunc("GET  /admin/metrics", config.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", config.handlerReset)
	return mux
}
