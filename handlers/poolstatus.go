package handlers

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"ck-pool-api/models"
)

func GetPoolStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := fmt.Sprintf("%s/logs/pool/pool.status", os.Getenv("POOL_BASE_PATH"))

		// Open the pool.status file
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Error opening pool.status file", http.StatusInternalServerError)
			log.Printf("Error opening pool.status: %v", err)
			return
		}
		defer file.Close()

		// Create a map to store the merged JSON data
		mergedData := make(map[string]interface{})
		reader := bufio.NewReader(file)

		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				http.Error(w, "Error reading pool.status", http.StatusInternalServerError)
				log.Printf("Error reading pool.status: %v", err)
				return
			}

			if len(line) > 0 {
				// Unmarshal each line into a temporary map
				var partialData map[string]interface{}
				if err := json.Unmarshal([]byte(line), &partialData); err != nil {
					log.Printf("Error parsing pool.status line: %v", err)
					continue // Skip invalid lines
				}

				// Merge partialData into mergedData
				for key, value := range partialData {
					mergedData[key] = value
				}
			}

			// Break at the end of the file
			if err == io.EOF {
				break
			}
		}

		// Convert the merged data map into a PoolStatus struct
		var poolStatus models.PoolStatus
		dataBytes, err := json.Marshal(mergedData)
		if err != nil {
			http.Error(w, "Error merging pool status data", http.StatusInternalServerError)
			log.Printf("Error merging pool status data: %v", err)
			return
		}

		// Unmarshal the merged data into the PoolStatus struct
		if err := json.Unmarshal(dataBytes, &poolStatus); err != nil {
			http.Error(w, "Error converting to PoolStatus struct", http.StatusInternalServerError)
			log.Printf("Error converting to PoolStatus struct: %v", err)
			return
		}

		// Return the merged pool status as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(poolStatus)
	}
}

// GetPoolHashratesHandler returns hashrate data for the pool
func GetPoolHashratesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT hashrate1m, hashrate5m, hashrate15m, hashrate1hr, hashrate6hr, hashrate1d, hashrate7d, saved_at 
			FROM pool_status 
			ORDER BY saved_at ASC`)

		if err != nil {
			http.Error(w, "Error fetching pool hashrates", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []map[string]interface{}

		for rows.Next() {
			var hashrate1m, hashrate5m, hashrate15m, hashrate1hr, hashrate6hr, hashrate1d, hashrate7d string
			var savedAt string

			if err := rows.Scan(&hashrate1m, &hashrate5m, &hashrate15m, &hashrate1hr, &hashrate6hr, &hashrate1d, &hashrate7d, &savedAt); err != nil {
				http.Error(w, "Error reading pool data", http.StatusInternalServerError)
				return
			}

			// Append the hashrate data and timestamp to the response
			data = append(data, map[string]interface{}{
				"hashrate1m":  hashrate1m,
				"hashrate5m":  hashrate5m,
				"hashrate15m": hashrate15m,
				"hashrate1hr": hashrate1hr,
				"hashrate6hr": hashrate6hr,
				"hashrate1d":  hashrate1d,
				"hashrate7d":  hashrate7d,
				"saved_at":    savedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
