package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hugobelem/chirpy/internal/auth"
)

func (config *apiConfig) handlerPolkaWebHook(
	w http.ResponseWriter,
	r *http.Request,
) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"could not resolve api key",
			err,
		)
		return
	}

	if apiKey != config.polkaKey {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"invalid api key",
			err,
		)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't covert ID to UUID type",
			err,
		)
		return
	}

	_, err = config.db.MarkUserAsChirpyRed(r.Context(), userUUID)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Couldn't update user",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
