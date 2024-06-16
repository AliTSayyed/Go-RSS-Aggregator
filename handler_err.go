package main

import "net/http"

// function to defeine a http reponse error handler in go
func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
