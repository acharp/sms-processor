package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Result is a response from API.
type Result struct {
	Message string `json:"message,omitempty"`
	// Status is a result status, can be either "ok" or "error".
	// The value is populated automatically based on the returning status code.
	Status string `json:"status"`
}

// JSONResponse responds with an application/json payload.
func JSONResponse(w http.ResponseWriter, msg string, statusCode int) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var result Result
	if statusCode >= 200 && statusCode < 300 {
		result = Result{Message: msg, Status: "ok"}
	} else {
		result = Result{Message: msg, Status: "error"}
	}

	content, err := json.Marshal(result)
	if err != nil {
		log.Printf("Marshalling failed: %s", err)
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(content)
	if err != nil {
		log.Printf("Writing failed: %s", err)
	}
}
