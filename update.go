package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hugobelem/chirpy/internal/auth"
	"github.com/hugobelem/chirpy/internal/database"
)

func (config *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"could not resolve token",
			err,
		)
		return
	}
	userID, err := auth.ValidateJWT(token, os.Getenv("SECRET"))
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"invalid token",
			err,
		)
		return
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't hash password",
			err,
		)
		return
	}

	user, err := config.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		log.Printf("Could't update user %s", err)
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could't update user",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        userID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
