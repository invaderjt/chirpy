package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/invaderjt/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirps(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, err, "error decoding JSON", 500)
		return
	}
	valid, cleaned_body := validateChirp(params.Body)
	if !valid {
		respondWithError(writer, err, "Chirp is too long", 400)
		return
	}

	newChirp := database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: params.UserID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), newChirp)
	if err != nil {
		respondWithError(writer, err, "could not submit chirp", 500)
	}

	successfulChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(writer, 201, successfulChirp)

}

func validateChirp(text string) (bool, string) {

	if len(text) > 140 {
		return false, text
	}
	return true, profanityCheck(text)
}

func profanityCheck(chirp string) string {
	profanityList := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(chirp, " ")
	for index, word := range words {
		for _, profanity := range profanityList {
			if strings.ToLower(word) == profanity {
				words[index] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
