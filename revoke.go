package main

import (
	"net/http"

	"github.com/hugobelem/chirpy/internal/auth"
)

func (config *apiConfig) handlerRevokeToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"could not resolve refresh token",
			err,
		)
		return
	}

	_, err = config.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"couldn't revoke session",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
