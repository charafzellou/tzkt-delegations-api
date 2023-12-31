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
		params := c.Queries()
		log.Printf("Params: %v\n", params)
		log.Printf("params[\"year\"]: %v\n", params["year"])
		error := c.JSON(streamDelegations(db, params["year"]))
		if error != nil {
			log.Fatalf("Error: %v\n", error)
		}
		var result Response
		result.Data = streamDelegations(db, params["year"])
		return c.JSON(result)
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

func streamDelegations(db *sql.DB, requested_year string) []DelegationApi {
	// Create a slice of DelegationApi
	delegations := []DelegationApi{}

	var rows *sql.Rows
	var err error
	// Query the database
	if requested_year == "" {
		rows, err = db.Query("SELECT timestamp, amount, sender, level FROM delegations")
	} else {
		rows, err = db.Query("SELECT timestamp, amount, sender, level FROM delegations WHERE EXTRACT(YEAR FROM timestamp) = $1", requested_year)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var delegation DelegationApi
		err := rows.Scan(&delegation.Timestamp, &delegation.Amount, &delegation.Delegator, &delegation.Block)
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
