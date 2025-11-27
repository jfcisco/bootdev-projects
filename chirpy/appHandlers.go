package main

import (
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Cast the anonymous func into an http.HandlerFunc type, which itself implements ServeHTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *chripyApp) mapAppHandlers() {
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	fileServer = c.cfg.middlewareMetricsInc(fileServer)
	c.mux.Handle("/app/", fileServer)
}
