package middleware

import "net/http"

func Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("API")
		if header != "ramsai" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
