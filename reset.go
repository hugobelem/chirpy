package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to reset the database: ",
			err,
		)
		return
	}

	type reset struct {
		Message string `json:"message"`
	}

	respondWithJSON(
		w,
		http.StatusOK,
		reset{
			Message: "Hits reset to 0 and database reset to initial state.",
		},
	)
}
