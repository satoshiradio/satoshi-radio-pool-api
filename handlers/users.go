package handlers

import (
	"ck-pool-api/models"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
)

var poolBasePath = os.Getenv("POOL_BASE_PATH")

func GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersDir := fmt.Sprintf("%s/logs/users", poolBasePath)

		// Open the users directory and list files
		files, err := ioutil.ReadDir(usersDir)
		if err != nil {
			http.Error(w, "Error reading users directory", http.StatusInternalServerError)
			log.Printf("Error reading users directory: %v", err)
			return
		}

		// Create a slice of usernames (file names)
		var users []string
		for _, file := range files {
			if !file.IsDir() {
				users = append(users, file.Name())
			}
		}

		// Return the list of users as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Regular expression for validating the username (as provided)
		regex := `^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^(bc1)[a-zA-HJ-NP-Z0-9]{39,59}$|^(bc1p)[a-zA-HJ-NP-Z0-9]{58}$`

		vars := mux.Vars(r)
		username := vars["username"]

		// Validate the username against the regex
		matched, err := regexp.MatchString(regex, username)
		if err != nil || !matched {
			http.Error(w, "Invalid username format", http.StatusBadRequest)
			return
		}
		filePath := fmt.Sprintf("%s/logs/users/%s", poolBasePath, username)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Open and read the user file
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading user file", http.StatusInternalServerError)
			log.Printf("Error reading user file %s: %v", username, err)
			return
		}

		// Unmarshal the file data into a User struct
		var user models.User
		if err := json.Unmarshal(fileBytes, &user); err != nil {
			http.Error(w, "Error parsing user file", http.StatusInternalServerError)
			log.Printf("Error parsing user file %s: %v", username, err)
			return
		}

		// Return the parsed user data as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// GetWorkerHashratesHandler returns hashrate data for a user's worker
func GetWorkerHashratesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		workername := vars["workername"]

		rows, err := db.Query(`SELECT hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, saved_at 
			FROM user_workers 
			WHERE username = $1 AND workername = $2 
			ORDER BY saved_at ASC`, username, workername)

		if err != nil {
			http.Error(w, "Error fetching worker hashrates", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []map[string]interface{}

		for rows.Next() {
			var hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d string
			var savedAt string

			if err := rows.Scan(&hashrate1m, &hashrate5m, &hashrate1hr, &hashrate1d, &hashrate7d, &savedAt); err != nil {
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}

			// Append the worker's hashrate data and timestamp to the response
			data = append(data, map[string]interface{}{
				"hashrate1m":  hashrate1m,
				"hashrate5m":  hashrate5m,
				"hashrate1hr": hashrate1hr,
				"hashrate1d":  hashrate1d,
				"hashrate7d":  hashrate7d,
				"saved_at":    savedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

// GetUserHashratesHandler returns hashrate data for a user
func GetUserHashratesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		rows, err := db.Query(`SELECT hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, saved_at 
			FROM users 
			WHERE username = $1 
			ORDER BY saved_at ASC`, username)

		if err != nil {
			println(err.Error())
			http.Error(w, "Error fetching hashrates", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []map[string]interface{}

		for rows.Next() {
			var hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d string
			var savedAt string

			if err := rows.Scan(&hashrate1m, &hashrate5m, &hashrate1hr, &hashrate1d, &hashrate7d, &savedAt); err != nil {
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}

			// Append the hashrate data and timestamp to the response
			data = append(data, map[string]interface{}{
				"hashrate1m":  hashrate1m,
				"hashrate5m":  hashrate5m,
				"hashrate1hr": hashrate1hr,
				"hashrate1d":  hashrate1d,
				"hashrate7d":  hashrate7d,
				"saved_at":    savedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
