package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/valbertoenoc/chirpy/internal/auth"
	"github.com/valbertoenoc/chirpy/internal/database"
)

// Parameters for creating a user
type createUserParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Return values for user creation
type createUserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// @Summary Create a user
// @Description Create a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body createUserParameters true "User parameters"
// @Success 201 {object} createUserResponse
// @Failure 400 {object} main.errorResponse
// @Failure 500 {object} main.errorResponse
// @Router /users [post]
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params createUserParameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error deconding request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "incorrect email or password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, createUserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}

// Parameters for updating a user
type updateUserParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response for user update
type updateUserResponse struct {
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ID          uuid.UUID `json:"id"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// @Summary Update a user
// @Description Update user email and password
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user body updateUserParameters true "User parameters"
// @Success 200 {object} updateUserResponse
// @Failure 400 {object} main.errorResponse
// @Failure 401 {object} main.errorResponse
// @Failure 500 {object} main.errorResponse
// @Router /users [put]
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no JWT found", err)
		return
	}

	validatedID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid", err)
		return
	}

	var params updateUserParameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             validatedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed at update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updateUserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}
