package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	// Define Postgres variables
	pg_host, pg_host_is_set         = os.LookupEnv("POSTGRES_HOST")
	pg_port, pg_port_is_set         = os.LookupEnv("POSTGRES_PORT")
	pg_user, pg_user_is_set         = os.LookupEnv("POSTGRES_USERNAME")
	pg_password, pg_password_is_set = os.LookupEnv("POSTGRES_PASSWORD")
	pg_db, pg_db_is_set             = os.LookupEnv("POSTGRES_DB")
)

func connectDatabase() *sql.DB {
	// Connect to the Postgres database
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg_user,
		pg_password,
		pg_host,
		pg_port,
		pg_db)
	fmt.Printf("postgresUrl: %s\n", postgresUrl)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initDatabase(db *sql.DB) {
	// Create the delegations table if it doesn't exist
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS delegations (hash TEXT PRIMARY KEY, level INT, timestamp TIMESTAMP, sender TEXT, delegate TEXT, amount FLOAT, status TEXT)")
	if err != nil {
		log.Fatal("[initDatabase] ", err)
	} else {
		log.Println("[initDatabase]  Table created")
	}
}

func getDelegationsCount() int {
	// Request delegations count from https://api.tzkt.io/v1/operations/drain_delegate/count
	response, err := http.Get("https://api.tzkt.io/v1/operations/delegations/count")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Print the HTTP Status Code and Status Name
	log.Println("[getDelegationsCount]  HTTP Response Status: ", response.StatusCode, http.StatusText(response.StatusCode))

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the response body to uint
	var delegationsCount int
	err = json.Unmarshal(body, &delegationsCount)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[getDelegationsCount]  DelegationsCount: ", delegationsCount)
	}

	return delegationsCount
}

func getDelegations() []Delegation {
	var delegations []Delegation
	// Request delegations from https://api.tzkt.io/v1/operations/delegations
	// for i := 0; i < getDelegationsCount(); i += 10000 {
	for i := 0; i < 1; i += 10000 {
		// url := fmt.Sprintf("https://api.tzkt.io/v1/operations/delegations?limit=10000&offset=%d", i)
		url := fmt.Sprintf("https://api.tzkt.io/v1/operations/delegations?limit=100&offset=%d", i)
		log.Printf("[getDelegations]  URL: %s\n", url)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// Print the HTTP Status Code and Status Name
		log.Println("[getDelegations]  HTTP Response Status:", response.StatusCode, http.StatusText(response.StatusCode))

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshal the JSON data
		var delegationsResponse []Delegation
		err = json.Unmarshal(body, &delegationsResponse)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("[getDelegations]  Delegations: ", delegationsResponse)
		}

		delegations = append(delegations, delegationsResponse...)

		// Sleep for 500 Milliseconds to avoid spamming the API
		// Tzkt.io has no rate limits currently, but it's better to be safe than sorry
		// This would be a good place to calculate a proper rate limit if we were using a different API
		time.Sleep(500 * time.Millisecond)
	}
	// fmt.Printf("[getDelegations] Delegations: %v\n", delegations)
	return delegations
}

func pushDelegations(delegations []Delegation, db *sql.DB) {
	fmt.Printf("[pushDelegations] Delegations Count : %d\n", len(delegations))
	// For each delegation in the delegations array
	for i := 0; i < len(delegations); i++ {
		// Check if the delegation already exists in the database
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM delegations WHERE hash=" + "'" + delegations[i].Hash + "'" + ")").Scan(&exists)
		if err != nil {
			log.Fatal(err)
		} else {
			// Insert Delegation into the database
			if !exists {
				_, err := db.Exec("INSERT INTO delegations (hash, level, timestamp, sender, delegate, amount, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					delegations[i].Hash,
					delegations[i].Level,
					delegations[i].Timestamp,
					delegations[i].Sender.Address,
					delegations[i].NewDelegate.Address,
					delegations[i].Amount,
					delegations[i].Status)
				if err != nil {
					log.Fatal(err)
				} else {
					log.Printf("[pushDelegations]  Delegation inserted: %d\n", i)
				}
			} else {
				log.Printf("[pushDelegations]  Delegation already exists: %v\n", delegations[i])
			}
		}
	}
}

func run() {
	// Verify that the required environment variables are set
	if !pg_user_is_set || !pg_host_is_set || !pg_port_is_set || !pg_password_is_set || !pg_db_is_set {
		fmt.Printf("[run] PG_HOST:    \"%s\"  =>  %t\n", pg_host, pg_host_is_set)
		fmt.Printf("[run] PG_PORT:    \"%s\"  =>  %t\n", pg_port, pg_port_is_set)
		fmt.Printf("[run] PG_USER:    \"%s\"  =>  %t\n", pg_user, pg_user_is_set)
		fmt.Printf("[run] PG_DB:      \"%s\"  =>   %t\n", pg_db, pg_db_is_set)
		log.Fatal("[run] Please set the required environment variables")
	}

	// Connect to the database
	db := connectDatabase()
	// Create necessary tables
	initDatabase(db)
	//  Get delegations from Tzkt.io
	delegations := getDelegations()
	// Push delegations to the database
	pushDelegations(delegations, db)

	// Close the database connection
	defer db.Close()
}

func main() {
	// Run the Indexer to get & push Delegations
	run()
}
