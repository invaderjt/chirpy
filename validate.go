package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error        error  `json:"error"`
		Valid        bool   `json:"valid"`
		Cleaned_Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, err, "error decoding JSON", 500)
		return
	}

	valid := false
	if len(params.Body) <= 140 {
		valid = true
	}

	if valid {
		body := profanityCheck(params.Body)
		respBody := returnVals{
			Valid:        valid,
			Cleaned_Body: body,
		}
		respondWithJSON(writer, 200, respBody)
		return
	}

	respondWithError(writer, nil, "Chirp is too long", 400)
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
