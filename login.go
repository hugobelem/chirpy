package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hugobelem/chirpy/internal/auth"
)

func (config *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	user, err := config.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println(err)
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}

	expiresIn := int64(3600)
	if params.ExpiresInSeconds > 0 {
		expiresIn = min(params.ExpiresInSeconds, 3600)
	}

	token, err := auth.MakeJWT(
		user.ID,
		os.Getenv("SECRET"),
		time.Duration(expiresIn),
	)
	if err != nil {
		log.Println(err)
		respondWithError(
			w,
			http.StatusInternalServerError,
			"an unexpected error accurred",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})
}
