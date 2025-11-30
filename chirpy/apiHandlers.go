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

		// Validate length
		if len(chirp.Body) > 140 {
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

		// Attempt to send success response
		res, err := json.Marshal(validateChirpResSuccess{true})
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

		err = rb.AddHeader("Content-Type", "application/json").
			Status(http.StatusOK).
			Bytes(res).
			WriteTo(w)
		logErr(err)
	})
}
