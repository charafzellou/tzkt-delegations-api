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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
)

var (
	config = fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "",
		AppName:       "",
	}
	recover_config = recover.Config{
		Next:              nil,
		EnableStackTrace:  false,
		StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
	}
	logger_config = logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "01-Jan-2000",
		TimeZone:   "Europe/Paris",
		Output:     setLogFile(),
	}
)

// Define a Delegation struct
type Delegation struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Level     int    `json:"level"`
	Timestamp string `json:"timestamp"`
	Block     string `json:"block"`
	Hash      string `json:"hash"`
	Counter   int    `json:"counter"`
	Initiator struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"initiator"`
	Sender struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"sender"`
	SenderCodeHash int `json:"senderCodeHash"`
	Nonce          int `json:"nonce"`
	GasLimit       int `json:"gasLimit"`
	GasUsed        int `json:"gasUsed"`
	StorageLimit   int `json:"storageLimit"`
	BakerFee       int `json:"bakerFee"`
	Amount         int `json:"amount"`
	PrevDelegate   struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"prevDelegate"`
	NewDelegate struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"newDelegate"`
	Status string `json:"status"`
	Errors []struct {
		Type string `json:"type"`
	} `json:"errors"`
	Quote struct {
		Btc int `json:"btc"`
		Eur int `json:"eur"`
		Usd int `json:"usd"`
		Cny int `json:"cny"`
		Jpy int `json:"jpy"`
		Krw int `json:"krw"`
		Eth int `json:"eth"`
		Gbp int `json:"gbp"`
	} `json:"quote"`
}

func setLogFile() *os.File {
	log_file, err := os.OpenFile("./logs/"+time.Now().String()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer log_file.Close()
	return log_file
}

func setApiServer() {
	app := fiber.New(config)
	app.Use(recover.New(recover_config))
	app.Use(logger.New(logger_config))

	// Create a "/xtz" endpoint group
	v1 := app.Group("/xtz")

	// Bind endpoints to handlers
	v1.Get("/delegations", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Setup static files
	app.Static("/", "./frontend/public")

	// Handle 404 : not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendFile("./frontend/public/404.html")
	})

	// Listen on port 3000
	log.Fatal(app.Listen(":3000"))
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
	for i := 0; i < getDelegationsCount(); i += 10000 {
		url := fmt.Sprintf("https://api.tzkt.io/v1/operations/delegations?limit=10000&offset=%d", i)
		log.Println("[getDelegations]  URL: ", url)
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
	return delegations
}

func pushDelegations() {
	// Connect to the Postgres database
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the array of delegations into the database
	delegations := getDelegations()
	for i := 0; i < len(delegations); i++ {
		// Check if the delegation already exists in the database
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM delegations WHERE hash=$1)", delegations[i].Hash).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		} else {
			// Insert Delegation into the database
			if !exists {
				_, err := db.Exec("INSERT INTO delegations (hash, level, timestamp, sender, delegate, quote, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					delegations[i].Hash,
					delegations[i].Level,
					delegations[i].Timestamp,
					delegations[i].Sender.Address,
					delegations[i].NewDelegate.Address,
					delegations[i].Quote.Usd,
					delegations[i].Status)
				if err != nil {
					log.Fatal(err)
				} else {
					log.Println("[pushDelegations]  Delegation inserted: ", delegations[i])
				}
			} else {
				log.Println("[pushDelegations]  Delegation already exists: ", delegations[i])
			}
		}
	}
}

func initDatabase() {
	// Connect to the Postgres database
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the delegations table if it doesn't exist
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS delegations (hash TEXT PRIMARY KEY, level INT, timestamp INT, sender TEXT, delegate TEXT, quote FLOAT, status TEXT)")
	// if err != nil {
	// 	log.Fatal("[initDatabase] ", err)
	// } else {
	// 	log.Println("[initDatabase]  Table created")
	// }
}

func main() {
	initDatabase()
	pushDelegations()
	setLogFile()
	setApiServer()
}
