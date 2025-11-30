package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type validateChirpReq struct {
	Body string `json:"body"`
}

type validateChirpResSuccess struct {
	CleanedBody string `json:"cleaned_body"`
}

type validateChirpResFailure struct {
	Error string `json:"error"`
}

// Helper func for a generic server error
func unexpectedErrorResponse(w http.ResponseWriter) {
	rb := NewResponseBuilder()
	err := rb.AddHeader("Content-Type", "text/plain; charset=utf-8").
		Status(http.StatusInternalServerError).
		Body("Something went wrong").
		WriteTo(w)
	logErr(err)
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
		rb := NewResponseBuilder()

		if err != nil {
			msg := "Something went wrong while decoding request"
			res, err := json.Marshal(validateChirpResFailure{msg})

			if err != nil {
				unexpectedErrorResponse(w)
				return
			}

			err = rb.AddHeader("Content-Type", "application/json").
				Status(http.StatusBadRequest).
				Bytes(res).
				WriteTo(w)
			logErr(err)
			return
		}

		fmt.Printf("[INFO] Received chirp: %s\n", chirp.Body)

		// Validate length
		if len(chirp.Body) > 140 {
			fmt.Println("[WARN] Chirp too long")
			msg := "Chirp is too long"
			res, err := json.Marshal(validateChirpResFailure{msg})

			if err != nil {
				unexpectedErrorResponse(w)
				return
			}

			err = rb.AddHeader("Content-Type", "application/json").
				Status(http.StatusBadRequest).
				Bytes(res).
				WriteTo(w)
			logErr(err)
			return
		}

		// Remove profane words from request body
		cleaned := removeProfaneWords(chirp.Body)

		// Attempt to send success response
		res, err := json.Marshal(validateChirpResSuccess{cleaned})
		if err != nil {
			msg := "Something went wrong while encoding response"
			res, err := json.Marshal(validateChirpResFailure{msg})

			if err != nil {
				unexpectedErrorResponse(w)
				return
			}

			err = rb.AddHeader("Content-Type", "application/json").
				Status(http.StatusInternalServerError).
				Bytes(res).
				WriteTo(w)
			logErr(err)
			return
		}

		fmt.Println("[INFO] Validate chirp success")

		err = rb.AddHeader("Content-Type", "application/json").
			Status(http.StatusOK).
			Bytes(res).
			WriteTo(w)
		logErr(err)
	})
}
