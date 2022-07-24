package middleware

import (
	"log"
	"net/http"
)

// TODO: allow auth bypass for local testing
func AuthVerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the API token from the header, and then verify the requestor
		handleAuth()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func handleAuth() {
	log.Println("handle auth here")
}
