package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jfcisco/boot-dev/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type chirpyAppCtx struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	mux            *http.ServeMux
}

func main() {
	// Load app configuration
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)

	// Set up database
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	app := &chirpyAppCtx{
		db:  dbQueries,
		mux: http.NewServeMux(),
	}

	// Map handlers
	app.mapAppHandlers()
	app.mapApiHandlers()
	app.mapAdminHandlers()

	// Serve app
	server := http.Server{
		Addr:    ":8080",
		Handler: app.mux,
	}

	fmt.Println("Listening on :8080...")
	fmt.Println("\nPages:")
	fmt.Println("http://localhost:8080/app")
	fmt.Println("http://localhost:8080/admin/metrics")
	log.Fatal(server.ListenAndServe())
}
