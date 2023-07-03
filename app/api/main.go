package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
)

var (
	// Define API port
	api_port, api_port_is_set = os.LookupEnv("API_PORT")

	// Define Postgres variables
	pg_host, pg_host_is_set         = os.LookupEnv("POSTGRES_HOST")
	pg_port, pg_port_is_set         = os.LookupEnv("POSTGRES_PORT")
	pg_user, pg_user_is_set         = os.LookupEnv("POSTGRES_USERNAME")
	pg_password, pg_password_is_set = os.LookupEnv("POSTGRES_PASSWORD")
	pg_db, pg_db_is_set             = os.LookupEnv("POSTGRES_DB")
)

func setLogFile() *os.File {
	log_file, err := os.OpenFile("./logs/"+time.Now().String()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	return log_file
}

func startApiServer(db *sql.DB) {
	app := fiber.New(config)
	app.Use(recover.New(recover_config))
	app.Use(logger.New(logger_config))

	// Create a "/xtz" endpoint group
	v1 := app.Group("/xtz")

	// Bind endpoints to handlers
	v1.Get("/delegations", func(c *fiber.Ctx) error {
		// return c.SendString("Hello, World 👋!")
		// read request body and unmarshal into Request struct
		// fmt.Printf("c.Body(): %v\n", c.Body())
		return c.JSON(streamDelegations(db))
	})

	// Setup static files
	app.Static("/", "./frontend/public")

	// Handle 404 : not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendFile("./frontend/public/404.html")
	})

	// Listen on port 3000
	listenerAddress := fmt.Sprintf(":%s", api_port)
	log.Printf("Listening on %s\n", listenerAddress)
	log.Fatal(app.Listen(listenerAddress))
}

func streamDelegations(db *sql.DB) []DelegationApi {
	// Create a slice of DelegationApi
	delegations := []DelegationApi{}

	// Query the database
	rows, err := db.Query("SELECT * FROM delegations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var delegation DelegationApi
		err := rows.Scan(&delegation.Hash, &delegation.Level, &delegation.Timestamp, &delegation.SenderAddress, &delegation.NewDelegateAddress, &delegation.Amount, &delegation.Status)
		if err != nil {
			log.Fatal(err)
		}
		delegations = append(delegations, delegation)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return delegations
}

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

func run() {
	// Verify that the required environment variables are set
	if !pg_user_is_set || !pg_host_is_set || !pg_port_is_set || !pg_password_is_set || !pg_db_is_set {
		log.Printf("[run]  PG_HOST: \"%s\" => %t\n", pg_host, pg_host_is_set)
		log.Printf("[run]  PG_PORT: \"%s\" => %t\n", pg_port, pg_port_is_set)
		log.Printf("[run]  PG_USER: \"%s\" => %t\n", pg_user, pg_user_is_set)
		log.Printf("[run]  PG_DB: \"%s\" => %t\n", pg_db, pg_db_is_set)
		log.Fatal("[run]  Please set the required environment variables")
	}

	// Set the API port to default as 3000 if not set
	if !api_port_is_set {
		api_port = "3000"
		log.Printf("[run]  API_PORT: \"%s\" => %t\n", api_port, api_port_is_set)
	}

	// Connect to the database
	db := connectDatabase()
	// Set the log file for the API server
	logFile := setLogFile()
	// Start the API server
	startApiServer(db)

	defer logFile.Close()
	defer db.Close()
}

func main() {
	// Run the application
	run()
}
