package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/suckoja/chirpy/internal/app"
	"github.com/suckoja/chirpy/internal/database"
	"github.com/suckoja/chirpy/internal/metrics"
)

func main() {
	// Prepare database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("DB_URL is required")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening db: %s\n", err)
		return
	}
	defer db.Close()

	dbQueries := database.New(db)

	// Prepare PageStat and Routing
	stats := &metrics.PageStats{}
	srv := app.NewServer(stats, dbQueries)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", srv.Routes()); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
