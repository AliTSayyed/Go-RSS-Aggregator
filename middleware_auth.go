package main

import (
	"fmt"
	"net/http"

	"github.com/alitsayyed/rssaggregator/internal/auth"
	"github.com/alitsayyed/rssaggregator/internal/database"
)

// Creat authHandler type to 'DRY' up the code. Any operation that needs a user's api key will go through this method.
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// but a reponse handler can only have the first 2 paramters, so use this method to convert a user to a http handler function
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// return a function (closure) that is a repsonse handler for a specific user
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKEY(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
