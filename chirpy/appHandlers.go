package main

import (
	"net/http"
)

func (c *chirpyAppCtx) middlewareMetricsInc(next http.Handler) http.Handler {
	// Cast the anonymous func into an http.HandlerFunc type, which itself implements ServeHTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *chirpyAppCtx) mapAppHandlers() {
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	fileServer = c.middlewareMetricsInc(fileServer)
	c.mux.Handle("/app/", fileServer)
}
