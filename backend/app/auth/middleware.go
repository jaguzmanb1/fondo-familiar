package auth

import (
	"fmt"
	"fondo-mod/data"
	"net/http"

	"github.com/gorilla/context"
)

// MiddlewareTokenValidationRol3 validates resquests tokens
func (h *Auth) MiddlewareTokenValidationRol3(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareTokenValidationRol3] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if r.Header["Authorization"] != nil {
			tv, us, err := h.validateToken(r.Header["Authorization"][0], "3")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.l.Error("[MiddlewareTokenValidationRol3] Error parsing or validating request token", "error", err, "endpoint", r.URL)
				data.ToJSON(&AuthError{Message: err.Error()}, w)

				return
			}

			if tv {
				context.Set(r, "us", us)
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		h.l.Info("[MiddlewareTokenValidationRol3] User request not autorized")
		fmt.Fprintf(w, "User request not autorized")
	})
}

// MiddlewareTokenValidationRol1 validates resquests tokens
func (h *Auth) MiddlewareTokenValidationRol1(next http.Handler) http.Handler {
	h.l.Info("[MiddlewareTokenValidationRol1] Handling validator middleware request")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		if r.Header["Authorization"] != nil {
			tv, us, err := h.validateToken(r.Header["Authorization"][0], "1")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				h.l.Error("[MiddlewareTokenValidationRol1] Error parsing or validating request token", "error", err, "endpoint", r.URL)

				return
			}

			if tv {
				context.Set(r, "us", us)
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		h.l.Info("[MiddlewareTokenValidationRol1] User request not autorized")
		fmt.Fprintf(w, "User request not autorized")
	})
}
