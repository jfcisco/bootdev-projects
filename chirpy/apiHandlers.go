package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type validateChirpReq struct {
	Body string `json:"body"`
}

type validateChirpResSuccess struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJsonResponse(rb *responseBuilder, statusCode int, payload any) {
	res, err := json.Marshal(payload)
	if err != nil {
		msg := "Something went wrong while encoding response"
		writeErrorResponse(rb, msg, http.StatusInternalServerError)
		return
	}

	err = rb.AddHeader("Content-Type", "application/json").
		Status(statusCode).
		Bytes(res).
		Write()
	logErr(err)
}

func writeErrorResponse(rb *responseBuilder, msg string, statusCode int) {
	res, err := json.Marshal(errorResponse{msg})
	if err != nil {
		err := rb.AddHeader("Content-Type", "text/plain; charset=utf-8").
			Status(http.StatusInternalServerError).
			Body("Something went wrong").
			Write()
		logErr(err)
	} else {
		err = rb.AddHeader("Content-Type", "application/json").
			Status(statusCode).
			Bytes(res).
			Write()
	}
	logErr(err)
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON from body
	var chirp validateChirpReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	rb := NewResponseBuilder(w)

	if err != nil {
		msg := "Something went wrong while decoding request"
		writeErrorResponse(rb, msg, http.StatusBadRequest)
		return
	}

	fmt.Printf("[INFO] Received chirp: %s\n", chirp.Body)

	// Validate length
	if len(chirp.Body) > 140 {
		fmt.Println("[WARN] Chirp too long")
		msg := "Chirp is too long"
		writeErrorResponse(rb, msg, http.StatusBadRequest)
		return
	}

	// Remove profane words from request body
	cleaned := removeProfaneWords(chirp.Body)
	fmt.Println("[INFO] Validate chirp success")
	writeJsonResponse(rb, http.StatusOK, validateChirpResSuccess{cleaned})
}

type createUserParams struct {
	Email string `json:"email"`
}

type createUserRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func CreateUserHandler(c *chirpyAppCtx) http.Handler {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Get email from request
			rb := NewResponseBuilder(w)
			decoder := json.NewDecoder(r.Body)

			var params createUserParams
			err := decoder.Decode(&params)
			if err != nil {
				writeErrorResponse(rb, "Unable to decode parameter from body", http.StatusBadRequest)
				return
			}

			if len(strings.TrimSpace(params.Email)) == 0 {
				writeErrorResponse(rb, "Cannot create user with no email", http.StatusBadRequest)
				return
			}

			newUser, err := c.db.CreateUser(r.Context(), params.Email)
			if err != nil {
				logErr(fmt.Errorf("error occurred while saving new user: %w", err))
				writeErrorResponse(rb, "An unexpected error occurred while saving user", http.StatusInternalServerError)
				return
			}

			fmt.Printf("[INFO] Successfully saved user with ID = %v\n", newUser.ID)

			res := createUserRes{
				ID:        newUser.ID,
				CreatedAt: newUser.CreatedAt,
				UpdatedAt: newUser.UpdatedAt,
				Email:     newUser.Email,
			}
			writeJsonResponse(rb, http.StatusCreated, res)
		},
	)
	return handler
}

func (c *chirpyAppCtx) mapApiHandlers() {
	c.mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	c.mux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)
	c.mux.Handle("POST /api/users", CreateUserHandler(c))
}
