package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hugobelem/chirpy/internal/auth"
	"github.com/hugobelem/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (config *apiConfig) handlerGetSingleChirp(
	w http.ResponseWriter,
	r *http.Request,
) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		log.Println("Chirp ID not provided")
		respondWithError(
			w,
			http.StatusBadRequest,
			"Chirp ID not provided",
			nil,
		)
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't covert ID to UUID type",
			err,
		)
		return
	}

	chirp, err := config.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (config *apiConfig) handlerRetrieveChirps(
	w http.ResponseWriter,
	r *http.Request,
) {
	chirps, err := config.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could't create chirp",
			err,
		)
	}

	listChirps := []Chirp{}
	for _, chirp := range chirps {
		listChirps = append(listChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, listChirps)

}

func (config *apiConfig) handlerCreateChirps(
	w http.ResponseWriter,
	r *http.Request,
) {
	type parameters struct {
		Body string `json:"body"`
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

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}

	chirp, err := config.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		log.Fatalf("Could't create chirp %s", err)
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could't create chirp",
			err,
		)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (config *apiConfig) handlerDeleteChirps(w http.ResponseWriter, r *http.Request) {
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

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		log.Println("Chirp ID not provided")
		respondWithError(
			w,
			http.StatusBadRequest,
			"Chirp ID not provided",
			nil,
		)
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't covert ID to UUID type",
			err,
		)
		return
	}

	chirp, err := config.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	if chirp.UserID != userID {
		respondWithError(
			w,
			http.StatusForbidden,
			"Coudn't delete chirp",
			nil,
		)
		return
	}

	_, err = config.db.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			"Coudn't delete chirp",
			nil,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	bardWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, bardWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
