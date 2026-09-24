package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsers(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, err, "error decoding JSON", 500)
		return
	}

	newUser, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(writer, err, "error creating new user", 500)
	}

	createdUser := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	respondWithJSON(writer, 201, createdUser)

}
