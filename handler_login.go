package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/valbertoenoc/chirpy/internal/auth"
	"github.com/valbertoenoc/chirpy/internal/database"
)

// Parameters for login
type loginParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Return values for login
type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

// Response for token refresh
type refreshResponse struct {
	Token string `json:"token"`
}

// @Summary Login a user
// @Description Login with email and password to get access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body loginParameters true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} main.errorResponse
// @Failure 401 {object} main.errorResponse
// @Failure 500 {object} main.errorResponse
// @Router /login [post]
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var params loginParameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error: %s", err), err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if ok, _ := auth.CheckPasswordHash(params.Password, user.HashedPassword); !ok {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expires_in := time.Hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, expires_in)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(expires_in),
	})

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed.Bool,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid Authorization header", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found or expired", err)
		return
	}

	// generate new token for corresponding user from the refresh token
	newToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed generating new token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid Authorization header", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "token revoked successfully")
}
