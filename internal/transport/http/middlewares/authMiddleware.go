package middlewares

import "net/http"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie := r.Header.Get("cookie")
		if cookie == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// TODO
	}
}
