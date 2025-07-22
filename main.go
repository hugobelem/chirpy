package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/hugobelem/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	port     = "8080"
	filepath = "/app/"
	dir      = "."
)

type apiConfig struct {
	fileHits atomic.Int32
	db       *database.Queries
	platform string
}

func main() {
	godotenv.Load()

	log.Printf(
		"Serving static files at %s from dir: %s on port: %s",
		filepath,
		dir,
		port,
	)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: setupRouter(),
	}

	log.Fatal(server.ListenAndServe())
}
