package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

/*
Purpose of this program is to build an RSS Aggregator.
RSS strands for Really Simple Syndication
and is an application that collects RSS feeds from multiple sources and displays them in one place.
An RSS feed is an online file that contains information about a website's published content.
It can include details like the content's full text or summary, publication date, author, and link.
RSS feed's data is written in XML.
An RSS Aggregator will collect different RSS feeds and add it to its database. Then it will automatically colellect
all the posts from those feeds and save them into the database. Then we can view the feeds and display them when / how we want.
*/

func main() {

	// download the godotenv dependency and add it to my go mod and store it locally in my vendor file.
	// this module allows me to pull information from my .env file automatically instead of manually having to write a command in the terminal every time I want to run my code (main).
	godotenv.Load(".env")

	// stores name of port from env file
	portString := os.Getenv("PORT")

	// if no port is found, use log fatal to exit the code immideately with error code 1 and a message
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	// create new router object to handle http requests. Use the go-chi dependency from github.
	router := chi.NewRouter()

	// create ability for users to send a request through a browser with different access types
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// create a new router for version 1, incase anything goes wrong it is nested in the version 1 path and doesnt affect other router paths
	v1Router := chi.NewRouter()
	// path inside v1 router, only works with get requests at the moment healthz is path name
	v1Router.Get("/healthz", handlerReadiness)
	// path inside v1 router for handling errors.
	v1Router.Get("/err", handlerErr)
	// nest the v1 router as its own path
	router.Mount("/v1", v1Router)

	// creating an instance of a http.Server struct called svr. Has 2 fields (Handler type and String)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	// listen and serve will block, when we get to this line our code stops and starts handling http requests.
	// if anything goes wrong in the process of handling those requests, then an error will be returned and the program will exit.
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)

}
