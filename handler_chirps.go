package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/valbertoenoc/chirpy/internal/auth"
	"github.com/valbertoenoc/chirpy/internal/database"
	"github.com/valbertoenoc/chirpy/internal/utils"
)

// Parameters for creating a chirp
type createChirpParameters struct {
	Body string `json:"body"`
}

// Response for chirp operations
type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
}

// @Summary Create a chirp
// @Description Create a new chirp with the given body
// @Tags chirps
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body createChirpParameters true "Chirp parameters"
// @Success 201 {object} chirpResponse
// @Failure 400 {object} main.errorResponse
// @Failure 401 {object} main.errorResponse
// @Failure 500 {object} main.errorResponse
// @Router /chirps [post]
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no JWT found", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid JWT token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params createChirpParameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to decode request body", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   utils.RedactProfanity(params.Body),
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		Token:     token,
		UserID:    userID,
	})
}

// @Summary List chirps
// @Description Get a list of chirps, optionally filtered by author ID
// @Tags chirps
// @Accept json
// @Produce json
// @Param author_id query string false "Filter by author ID"
// @Success 200 {array} database.Chirp
// @Failure 500 {object} main.errorResponse
// @Router /chirps [get]
func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	authorIDString := r.URL.Query().Get("author_id")
	authorID, err := uuid.Parse(r.URL.Query().Get("author_id"))
	if authorIDString != "" && err == nil {
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.ListChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// sort := r.URL.Query().Get("sort")
	// if sort == "" || sort == "asc" {
	// 	chirps.
	// }

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(r.PathValue("id"))
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := uuid.MustParse(r.PathValue("id"))

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate token", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not find chirp", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "user not allowed to delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
