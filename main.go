package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/brentjolicoeur/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	fSrv := http.FileServer(http.Dir(filepathRoot))
	fileserverHandler := http.StripPrefix("/app", fSrv)

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s.\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}
