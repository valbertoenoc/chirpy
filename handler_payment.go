package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/valbertoenoc/chirpy/internal/auth"
)

// Parameters for webhook
type webhookParameters struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

// @Summary Upgrade user to Chirpy Red
// @Description Handle webhook from Polka to upgrade user to Chirpy Red
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param body body webhookParameters true "Webhook parameters"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} main.errorResponse
// @Failure 401 {object} main.errorResponse
// @Failure 404 {object} main.errorResponse
// @Router /api/polka/webhooks [post]
func (cfg *apiConfig) handlerUpgradeToRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "API Key missing", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	var params webhookParameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "unsupported event", nil)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not find user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
