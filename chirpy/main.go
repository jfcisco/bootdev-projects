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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Cast the anonymous func into an http.HandlerFunc type, which itself implements ServeHTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := &apiConfig{}

	mux := http.NewServeMux()

	appFs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	appFs = cfg.middlewareMetricsInc(appFs)
	mux.Handle("/app/", appFs)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		doc := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin!</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, cfg.fileserverHits.Load())
		w.Write([]byte(doc))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Listening on :8080...")
	fmt.Println("\nPages:")
	fmt.Println("http://localhost:8080/app")
	fmt.Println("http://localhost:8080/admin/metrics")
	log.Fatal(server.ListenAndServe())
}
