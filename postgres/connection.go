package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	dsn := "host=aws-0-ap-southeast-1.pooler.supabase.com port=6543 user=postgres.vfhquqdbosumijamcqvo password=ratatouille@hello! dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}
	log.Println("Connected to the database successfully.")
	return db
}
