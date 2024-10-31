package db

import (
	"bufio"
	"ck-pool-api/models"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func InitDB() (*sql.DB, error) {
	connStr := getPostgresConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	createTables(db)

	return db, nil
}

// getPostgresConnStr constructs the PostgreSQL connection string from environment variables
func getPostgresConnStr() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// Construct PostgreSQL connection string
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func createTables(db *sql.DB) {
	// Pool status table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pool_status (
		id SERIAL PRIMARY KEY,
		runtime BIGINT,
		lastupdate BIGINT,
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
		accepted BIGINT, 
		rejected BIGINT, 
		bestshare BIGINT, 
		sps1m REAL, 
		sps5m REAL, 
		sps15m REAL, 
		sps1h REAL,
		saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal(err)
	}

	// Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT, 
		hashrate1m TEXT, 
		hashrate5m TEXT, 
		hashrate1hr TEXT, 
		hashrate1d TEXT, 
		hashrate7d TEXT, 
		lastshare BIGINT, 
		workers INTEGER, 
		shares BIGINT, 
		bestshare REAL, 
		bestever BIGINT, 
		authorised BIGINT,
		saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal(err)
	}

	// User workers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_workers (
		id SERIAL PRIMARY KEY,
		username TEXT, 
		workername TEXT,
		hashrate1m TEXT,
		hashrate5m TEXT,
		hashrate1hr TEXT,
		hashrate1d TEXT,
		hashrate7d TEXT,
		lastshare BIGINT,
		shares BIGINT,
		bestshare REAL,
		bestever BIGINT,
		saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	// Insert pool status into the database
	_, err = db.Exec(`INSERT INTO pool_status (runtime, lastupdate, users, workers, idle, disconnected, hashrate1m, hashrate5m, hashrate15m, hashrate1hr, hashrate6hr, hashrate1d, hashrate7d, diff, accepted, rejected, bestshare, sps1m, sps5m, sps15m, sps1h) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`,
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

			// Insert user data into the users table
			_, err = db.Exec(`INSERT INTO users (username, hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, lastshare, workers, shares, bestshare, bestever, authorised) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
				username, user.Hashrate1m, user.Hashrate5m, user.Hashrate1hr, user.Hashrate1d, user.Hashrate7d, user.LastShare, user.Workers, user.Shares, user.BestShare, user.BestEver, user.Authorised)

			if err != nil {
				log.Printf("Error inserting user data for %s into database: %v", username, err)
			}

			// Now insert each worker's data into the user_workers table
			for _, worker := range user.Worker {
				_, err = db.Exec(`INSERT INTO user_workers (username, workername, hashrate1m, hashrate5m, hashrate1hr, hashrate1d, hashrate7d, lastshare, shares, bestshare, bestever) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
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
