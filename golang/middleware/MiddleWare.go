package middleware

import (
	"log"
	"net/http"
)

const methodNotAllowed = "Method not allowed"

func AllowMethod(method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				log.Printf("Method %s not allowed for %s", r.Method, r.URL.Path)
				http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
