package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type chripyApp struct {
	cfg *apiConfig
	mux *http.ServeMux
}

func main() {
	app := &chripyApp{
		cfg: &apiConfig{},
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
