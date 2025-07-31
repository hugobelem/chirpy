package main

import (
	"net/http"

	"github.com/hugobelem/chirpy/internal/auth"
)

func (config *apiConfig) handlerRevokeToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	refreshTokenParam, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"could not resolve refresh token",
			err,
		)
		return
	}

	_, err = config.db.RevokeRefreshToken(r.Context(), refreshTokenParam)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"refresh token not found",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}