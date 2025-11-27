package main

import (
	"fmt"
	"net/http"
)

func (c *chripyApp) mapAdminHandlers() {
	c.mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		c.cfg.fileserverHits.Store(0)
	})

	c.mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		doc := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin!</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, c.cfg.fileserverHits.Load())
		w.Write([]byte(doc))
	})
}
