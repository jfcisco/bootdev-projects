package main

import (
	"encoding/json"
	"net/http"
)

type validateChirpReq struct {
	Body string `json:"body"`
}

type validateChirpResSuccess struct {
	Valid bool `json:"valid"`
}

type validateChirpResFailure struct {
	Error string `json:"error"`
}

func (c *chripyApp) mapApiHandlers() {
	c.mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	c.mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		// Decode JSON from body
		var chirp validateChirpReq
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&chirp)

		if err != nil {
			msg := "Something went wrong while decoding request"
			res, err := json.Marshal(validateChirpResFailure{msg})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something went wrong"))
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res)
			}
			return
		}

		// Validate length
		if len(chirp.Body) > 140 {
			msg := "Chirp is too long"
			res, err := json.Marshal(validateChirpResFailure{msg})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something went wrong"))
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		res, _ := json.Marshal(validateChirpResSuccess{true})
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})
}
