package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"ck-pool-api/db"
	"ck-pool-api/handlers"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}
}
func main() {
	// Initialize the SQLite database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Start goroutine for saving data every 5 minutes
	go func() {
		for {
			db.StorePoolStatus(database, "ckpool/logs/pool/pool.status")
			db.StoreUserFiles(database, "ckpool/logs/users")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Set up the API routes
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/pool", handlers.GetPoolStatusHandler()).Methods("GET")
	router.HandleFunc("/api/v1/pool/hashrates", handlers.GetPoolHashratesHandler(database)).Methods("GET")
	//router.HandleFunc("/api/v1/users", handlers.GetUsersHandler()).Methods("GET") // temporary disabled for privacy reasons
	router.HandleFunc("/api/v1/users/{username}", handlers.GetUserHandler()).Methods("GET")
	router.HandleFunc("/api/v1/users/{username}/hashrates", handlers.GetUserHashratesHandler(database)).Methods("GET")
	router.HandleFunc("/api/v1/users/{username}/workers/{workername}/hashrates", handlers.GetWorkerHashratesHandler(database)).Methods("GET")
	// Enable CORS for all routes
	headersOk := gohandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gohandlers.AllowedOrigins([]string{"*"})
	methodsOk := gohandlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	// Start the server with CORS enabled
	fmt.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", gohandlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
