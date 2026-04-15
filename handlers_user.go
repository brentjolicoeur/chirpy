package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brentjolicoeur/chirpy/internal/auth"
	"github.com/brentjolicoeur/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal response", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	userResponse := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, userResponse)
}

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal response", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token", err)
		return
	}
	refreshString := auth.MakeRefreshToken()

	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	userResponse := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't retrieve token", err)
		return
	}
	verifiedID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal response", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}
	updatedUserParams := database.UpdateUserParams{
		ID:             verifiedID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.UpdateUser(r.Context(), updatedUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update user", err)
		return
	}

	userResponse := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't retrieve apiKey", err)
		return
	}
	if apiKey != os.Getenv("POLKA_KEY") {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	defer r.Body.Close()

	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request", err)
		return
	}
	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal response", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid userID", err)
	}

	err = cfg.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
