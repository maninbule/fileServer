package middleware

import (
	"fmt"
	"net/http"
)

// Middleware to limit the size of POST requests
func LimitRequestBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("enter LimitRequestBodyMiddleware")
		const maxUploadSize = 22<<20 + 512 // 12.5MB
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		next.ServeHTTP(w, r)
	})
}
