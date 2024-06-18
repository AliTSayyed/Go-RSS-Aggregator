package main

import "net/http"

// function to defeine a http reponse handler
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
