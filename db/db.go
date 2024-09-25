package db

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"ck-pool-api/models"
	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./pooldata.db")
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	createTables(db)

	return db, nil
}

func createTables(db *sql.DB) {
	// Pool status table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pool_status (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		runtime INTEGER, 
		lastupdate INTEGER, 
		users INTEGER, 
		workers INTEGER, 
		idle INTEGER, 
		disconnected INTEGER, 
		hashrate1m TEXT, 
		hashrate5m TEXT, 
		hashrate15m TEXT, 
		hashrate1hr TEXT, 
		hashrate6hr TEXT, 
		hashrate1d TEXT, 
		hashrate7d TEXT, 
		diff REAL, 
		accepted INTEGER, 
		rejected INTEGER, 
		bestshare INTEGER, 
		sps1m REAL, 
		sps5m REAL, 
		sps15m REAL, 
		sps1h REAL,
    saved_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal(err)
	}

	// Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT, 
		hashrate1m TEXT, 
		hashrate5m TEXT, 
		hashrate1hr TEXT, 
		hashrate1d TEXT, 
		hashrate7d TEXT, 
		lastshare INTEGER, 
		workers INTEGER, 
		shares INTEGER, 
		bestshare REAL, 
		bestever INTEGER, 
		authorised INTEGER,
    saved_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_workers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,  -- Foreign key for user identification
    workername TEXT,
    hashrate1m TEXT,
    hashrate5m TEXT,
    hashrate1hr TEXT,
    hashrate1d TEXT,
    hashrate7d TEXT,
    lastshare INTEGER,
    shares INTEGER,
    bestshare REAL,
    bestever REAL,
    saved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(username) REFERENCES users(username)
);`)

	if err != nil {
		log.Fatal(err)
	}
}

func StorePoolStatus(db *sql.DB, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening pool.status: %v", err)
		return
	}
	defer file.Close()

	var poolStatus models.PoolStatus
	scanner := bufio.NewScanner(file)

	// Read each line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &poolStatus); err != nil {
			log.Printf("Error parsing pool.status: %v", err)
			continue // Skip lines that can't be parsed
		}
	}

	_, err = db.Exec(`INSERT INTO pool_status (runtime,lastupdate, users, workers, idle, disconnected, hashrate1m, hashrate5m, hashrate15m, hashrate1hr, hashrate6hr, hashrate1d, hashrate7d, diff, accepted, rejected, bestshare, sps1m, sps5m, sps15m, sps1h) 
	VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		poolStatus.Runtime, poolStatus.LastUpdate, poolStatus.Users, poolStatus.Workers, poolStatus.Idle, poolStatus.Disconnected,
		poolStatus.Hashrate1m, poolStatus.Hashrate5m, poolStatus.Hashrate15m, poolStatus.Hashrate1hr, poolStatus.Hashrate6hr, poolStatus.Hashrate1d, poolStatus.Hashrate7d,
		poolStatus.Diff, poolStatus.Accepted, poolStatus.Rejected, poolStatus.BestShare, poolStatus.SPS1m, poolStatus.SPS5m, poolStatus.SPS15m, poolStatus.SPS1h)
	if err != nil {
		log.Printf("Error inserting pool status into database: %v", err)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading pool.status: %v", err)
	}
}

func StoreUserFiles(db *sql.DB, usersDir string) {
	err := filepath.Walk(usersDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			username := info.Name() // Make sure you pass the username
			fileBytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading user file %s: %v", username, err)
				return nil
			}

			var user models.User
			if err := json.Unmarshal(fileBytes, &user); err != nil {
				log.Printf("Error parsing user file %s: %v", username, err)
				return nil
			}

			// Insert into database (ensure 12 values are passed)
			_, err = db.Exec(`INSERT INTO users (username, hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, lastshare, workers, shares, bestshare, bestever, authorised) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				username, user.Hashrate1m, user.Hashrate5m, user.Hashrate1hr, user.Hashrate1d, user.Hashrate7d, user.LastShare, user.Workers, user.Shares, user.BestShare, user.BestEver, user.Authorised)

			if err != nil {
				log.Printf("Error inserting user data for %s into database: %v", username, err)
			}
			// Now insert each worker's data into the user_workers table
			for _, worker := range user.Worker {
				_, err = db.Exec(`INSERT INTO user_workers (username, workername, hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, lastshare, shares, bestshare, bestever) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					username, worker.WorkerName, worker.Hashrate1m, worker.Hashrate5m, worker.Hashrate1hr, worker.Hashrate1d, worker.Hashrate7d, worker.LastShare, worker.Shares, worker.BestShare, worker.BestEver)

				if err != nil {
					log.Printf("Error inserting worker data for %s into database: %v", worker.WorkerName, err)
					return nil
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking users directory: %v", err)
	}
}
