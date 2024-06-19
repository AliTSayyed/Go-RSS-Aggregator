package main

// how we will display data/errors in json
import (
	"encoding/json"
	"log"
	"net/http"
)

// this function will deal with how to display various error messages
func respondWithError(w http.ResponseWriter, code int, msg string) {
	//any 500 error code means there is a bug on our end.
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	// create a formater error, uses a json tag
	type errResponse struct {
		Error string `json:"error"`
	}
	// responding with error is a responding with json function with a errResponse type at the end
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

// this function will deal with how to display various data and out put in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// payload is the data we are passing through
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marhsal JSON reponse: %v", payload)
		w.WriteHeader(500)
		return
	}
	// create a header to show we are repsonding with json data
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
