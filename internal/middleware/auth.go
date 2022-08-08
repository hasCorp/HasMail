package middleware

import (
	"net/http"
)

const HeaderAuthorization = "Authorization"

type LocalAuthMiddleware struct {
	AllowedToken string
}

func (m *LocalAuthMiddleware) LocalAuthVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the API token from the header, and then verify that the
		// token matches the allowed local token
		if !m.compareToLocalToken(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *LocalAuthMiddleware) compareToLocalToken(r *http.Request) bool {
	return m.AllowedToken == readAuthToken(r)
}

func AuthVerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the API token from the header, and then verify the requestor
		if !validateToken(readAuthToken(r)) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func readAuthToken(r *http.Request) string {
	token, ok := r.Header[HeaderAuthorization]
	if !ok || len(token) < 1 {
		return ""
	}
	return token[0]
}

func validateToken(token string) bool {
	return true // TODO: finish this when hasAuth is done
}
