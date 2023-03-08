package middlewares

import (
	"net/http"

	"github.com/projeto-BackEnd/security"
)

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		if r.URL.String() == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		errValidToken := security.ValidToken(r)
		if errValidToken != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(errValidToken.Error()))
			return
		}
		next.ServeHTTP(w, r)
	})
}
