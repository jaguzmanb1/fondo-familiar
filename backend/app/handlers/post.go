package handlers

import (
	"fondo-mod/data"
	"net/http"

	"github.com/gorilla/context"
)

//CreateAporte handles the request to create an aporte in the database
func (h *UsersHandler) CreateAporte(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	var ap = (context.Get(r, "ap")).(*data.Aporte)

	h.l.Info("[CreateAporte] Creating new aporte to user", "user", us)
	err := h.UserService.CreateAporte(ap.IDUsuario, ap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
}

//CreatePago handles the request to create a pago in the database
func (h *UsersHandler) CreatePago(w http.ResponseWriter, r *http.Request) {
	var p = (context.Get(r, "p")).(*data.Pago)

	h.l.Info("[CreateAporte] Creating new pago to credit", "credit", p)
	err := h.UserService.CreatePago(p)
	err = h.UserService.CreatePagoInteres(p)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
}

//CreateCredito handles the request to create a credito in the database
func (h *UsersHandler) CreateCredito(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	var cr = (context.Get(r, "cr")).(*data.Credito)

	h.l.Info("[CreateCredito] Creating new aporte to user", "user", us)
	err := h.UserService.CreateCredito(cr.IDUsuario, cr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
}
