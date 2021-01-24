package handlers

import (
	"fondo-mod/data"
	"net/http"

	"github.com/gorilla/context"
)

//MiddlewareValidateDescuento  verificacion para los request
func (h *UsersHandler) MiddlewareValidateDescuento(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		descuento := &data.PostDescuento{}

		err := data.FromJSON(descuento, r.Body)
		if err != nil {
			h.l.Error("deserializing descuento", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		h.l.Debug("Serialized aporte", "descuento", descuento)
		errs := h.v.Validate(descuento)
		if len(errs) != 0 {
			h.l.Error("validating descuento", "errors:", errs)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		context.Set(r, "d", descuento)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

//MiddlewareValidateAporte  verificacion para los request
func (h *UsersHandler) MiddlewareValidateAporte(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		aporte := &data.Aporte{}

		err := data.FromJSON(aporte, r.Body)
		if err != nil {
			h.l.Error("deserializing aporte", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		h.l.Debug("Serialized aporte", "aporte", aporte)
		errs := h.v.Validate(aporte)
		if len(errs) != 0 {
			h.l.Error("validating aporte", "errors:", errs)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		context.Set(r, "ap", aporte)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

//MiddlewareValidatePago  verificacion para los request
func (h *UsersHandler) MiddlewareValidatePago(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		credito := &data.Pago{}

		err := data.FromJSON(credito, r.Body)
		if err != nil {
			h.l.Error("[MiddlewareValidatePago] Deserializing pago", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		h.l.Debug("[MiddlewareValidatePago] Serialized credito", "pago", credito)
		errs := h.v.Validate(credito)
		if len(errs) != 0 {
			h.l.Error("[MiddlewareValidatePago] Validating pago", "errors:", errs)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		context.Set(r, "p", credito)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

//MiddlewareValidateCredito  verificacion para los request
func (h *UsersHandler) MiddlewareValidateCredito(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		credito := &data.Credito{}

		err := data.FromJSON(credito, r.Body)
		if err != nil {
			h.l.Error("[MiddlewareValidateCredito] Deserializing credito", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		h.l.Debug("[MiddlewareValidateCredito] Serialized credito", "credito", credito)
		errs := h.v.Validate(credito)
		if len(errs) != 0 {
			h.l.Error("[MiddlewareValidateCredito] Validating credito", "errors:", errs)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		context.Set(r, "cr", credito)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

//MiddlewareCheckUserIDCall verifies that the id sent from the user is the same as the speciefied on the token
func (h *UsersHandler) MiddlewareCheckUserIDCall(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		aporte := &data.Aporte{}
		var us = (context.Get(r, "us")).(data.User)
		id := getID(r)

		if us.Rol != 1 {
			if us.Rol != id {
				h.l.Error("User trying to access data from another user", "User Origin", us, "id", id)
				rw.WriteHeader(http.StatusUnauthorized)
				data.ToJSON(&ValidationError{Messages: []string{"User trying to access data from another user"}}, rw)

				return
			}
		}

		context.Set(r, "ap", aporte)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
