package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/hugobelem/chirpy/internal/auth"
)

func (config *apiConfig) handlerRefreshToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	type response struct{
		Token string `json:"token"`
	}

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

	userID, err := config.db.GetUserFromRefeshToken(
		r.Context(),
		refreshTokenParam,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"token does not exist or expired",
			err,
		)
		return
	}

	refreshToken, _ := config.db.GetRefreshToken(r.Context(), refreshTokenParam)

	notRevoked := sql.NullTime{}
	if refreshToken.RevokedAt != notRevoked {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(userID, os.Getenv("SECRET"))
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"could not generate new access token",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}