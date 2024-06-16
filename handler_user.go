package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alitsayyed/rssaggregator/internal/database"
	"github.com/google/uuid"
)

// Method to defeine a http reponse handler for an apiConfig struct. Gives handler access to the databse.
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// handler needs to take as input a json body
	type parameters struct {
		Name string `json:"name"`
	}

	// need to parse the request body into the paramter struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	// now we can create a user from the sqlc generated go code
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user : %v", err))
		return
	}

	// if user is created this will be the final json repsonse (use our template for user not sqlc)
	respondWithJSON(w, 200, databaseUserToUser(user))
}
