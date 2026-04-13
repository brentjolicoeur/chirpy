package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func cleanChirpText(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle":
			words[i] = "****"
			continue
		case "sharbert":
			words[i] = "****"
			continue
		case "fornax":
			words[i] = "****"
			continue
		default:
			continue
		}
	}

	return strings.Join(words, " ")
}
