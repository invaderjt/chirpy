package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(writer http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		respondWithError(writer, err, "could not fetch chirps", 500)
		return
	}

	var chirpFeed []Chirp
	for _, chirp := range chirps {
		nextChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpFeed = append(chirpFeed, nextChirp)
	}
	respondWithJSON(writer, 200, chirpFeed)

}

func (cfg *apiConfig) handlerGetChirp(writer http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, err, "Invalid chirp ID", http.StatusBadRequest)
		return
	}
	chirp, err := cfg.dbQueries.GetSingleChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(writer, err, "chirp not found", 404)
		return
	}

	requestedChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(writer, 200, requestedChirp)
}
