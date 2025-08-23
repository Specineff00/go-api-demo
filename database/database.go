package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	// Creates a new SQLite database file
	DB, err = sql.Open("sqlite3", "./database/users.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database successfully corrected")
}

func CreateTables() {
	// Read schema from file
	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Users tables successfully created")
}

func SeedData() {
	seedSQL, err := os.ReadFile("./database/seed.sql")

	if err != nil {
		log.Printf("Warning! Could not read seed file %v", err)
		return
	}

	// Execute the seed SQL
	_, err = DB.Exec(string(seedSQL))
	if err != nil {
		log.Printf("Warning: Could not seed data %v", err)
		return
	}

	log.Println("Database seeded with sample data")
}
