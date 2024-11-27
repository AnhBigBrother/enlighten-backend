package cfg

import (
	"database/sql"
	"enlighten-backend/internal/database"
	"log"
	"os"
)

var DBQueries *database.Queries

func init() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	DBQueries = database.New(conn)
}
