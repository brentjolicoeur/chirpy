package main

import (
	"net/http"
	"time"

	"github.com/brentjolicoeur/chirpy/internal/auth"
)

func (cfg *apiConfig) apiRefreshHandler(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing token in header", err)
		return
	}
	token, err := cfg.db.GetToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token doesn't exit", err)
		return
	}
	if token.ExpiresAt.Before(time.Now()) || token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token is invalid", nil)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving user", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating JWT", err)
		return
	}
	type responseBody struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, responseBody{
		Token: accessToken,
	})
}

func (cfg *apiConfig) apiRevokeHandler(w http.ResponseWriter, r *http.Request) {

}
