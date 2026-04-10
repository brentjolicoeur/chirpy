package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) adminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1 <p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal response", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, responseBody{
		Valid: true,
	})
}
