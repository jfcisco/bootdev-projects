package main

import (
	"fmt"
	"net/http"
)

func logErr(err error) {
	if err != nil {
		fmt.Printf("[ERR]: %v\n", err)
	}
}

type responseBuilder struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func NewResponseBuilder() *responseBuilder {
	return &responseBuilder{}
}

func (r *responseBuilder) Status(statusCode int) *responseBuilder {
	r.statusCode = statusCode
	return r
}

func (r *responseBuilder) Body(body string) *responseBuilder {
	r.body = []byte(body)
	return r
}

func (r *responseBuilder) Bytes(body []byte) *responseBuilder {
	r.body = body
	return r
}

func (r *responseBuilder) AddHeader(key string, value string) *responseBuilder {
	if r.headers == nil {
		r.headers = make(http.Header)
	}

	r.headers.Add(key, value)
	return r
}

func (r *responseBuilder) WriteTo(responseWriter http.ResponseWriter) error {
	if r.statusCode > 0 {
		responseWriter.WriteHeader(r.statusCode)
	}

	if r.headers != nil {
		for key, values := range r.headers {
			for _, value := range values {
				responseWriter.Header().Add(key, value)
			}
		}
	}

	if len(r.body) > 0 {
		_, err := responseWriter.Write(r.body)
		if err != nil {
			return err
		}
	}
	return nil
}
