package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hugobelem/chirpy/internal/auth"
	"github.com/hugobelem/chirpy/internal/database"
)

func (config *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	token, err := auth.MakeJWT(
		user.ID,
		os.Getenv("SECRET"),
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

	refreshToken, _ := auth.MakeRefreshToken()
	persistRefreshToken, err := config.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60),
		})
	if err != nil {
		log.Println(err)
		respondWithError(
			w,
			http.StatusInternalServerError,
			"could not persist refresh token",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: persistRefreshToken.Token,
		},
	})
}
