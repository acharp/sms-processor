package main

import (
	"net/http"
)

// authAPIKey does a very basic authorisation check from the API Key.
func authAPIKey(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("X-API-Key")
		if key != "BigBird" {
			JSONResponse(w, "Unauthorised key", http.StatusUnauthorized)
			return
		}
		handler(w, req)
	}
}
