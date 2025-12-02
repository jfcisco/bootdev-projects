package main

import (
	"fmt"
	"net/http"
	"os"
)

func (c *chirpyAppCtx) mapAdminHandlers() {
	c.mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		currentPlatform := os.Getenv("PLATFORM")
		if currentPlatform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		c.fileserverHits.Store(0)

		// Delete all users
		err := c.db.DeleteAllUsers(r.Context())
		if err != nil {
			logErr(err)
		}
	})

	c.mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		doc := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin!</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, c.fileserverHits.Load())
		w.Write([]byte(doc))
	})
}
