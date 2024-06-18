package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alitsayyed/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

// struct that holds a connection to a database
type apiConfig struct {
	DB *database.Queries
}

func main() {

	// download the godotenv dependency and add it to my go mod and store it locally in my vendor file.
	// this module allows me to pull information from my .env file automatically instead of manually having to write a command in the terminal every time I want to run my code (main).
	godotenv.Load(".env")

	// stores name of port from env file
	// if no port is found, use log fatal to exit the code immideately with error code 1 and a message
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	// stores name of the Database url from env file
	// if no DB url is found, use log fatal to exit the code immideately with error code 1 and a message
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// store a connection to  a database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant' connect to databse:", err)
	}

	// need to convert the sql.DB to a DB for the apiConfig struct so we use the database.New method
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	// call the scraper to create 10 concurrent go routines to fetch and update posts
	go startScraping(db, 10, time.Minute)

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
	// create handler that has access to the database
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	// create handler to get unique users and feeds from the database, uses middle ware function for user authentification
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	// delete method has a specified feed id to delete in its path
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
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
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)

}
