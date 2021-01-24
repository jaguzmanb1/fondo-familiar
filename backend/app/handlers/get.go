package handlers

import (
	"fondo-mod/data"
	"net/http"

	"github.com/gorilla/context"
)

// GetAllAportes returns all aportes in the fondo
func (h *UsersHandler) GetAllAportes(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)

	h.l.Info("[GetAllAportes] Recieving call to get all aportes from", "user", us)
	aportes, err := h.UserService.GetAllAportes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&aportes, w)
}

// GetAllCreditos returns all aportes in the fondo
func (h *UsersHandler) GetAllCreditos(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)

	h.l.Info("[GetAllCreditos] Recieving call to get all creditos from", "user", us)
	creditos, err := h.UserService.GetAllCreditos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&creditos, w)
}

// GetReporteGeneral returns all aportes in the fondo
func (h *UsersHandler) GetReporteGeneral(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)

	h.l.Info("[GetReporteGeneral] Recieving call to get a general report from ", "user", us)
	reporte, err := h.UserService.GetReporteGeneral()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&reporte, w)
}

// GetAllCreditosByUserID returns all creditos in the fondo given a user ID
func (h *UsersHandler) GetAllCreditosByUserID(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	id := getID(r)

	h.l.Info("[GetAllCreditosByUserID] Recieving request to get all credits from", "user", us)
	creditos, err := h.UserService.GetAllCreditosByUserID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&creditos, w)
}

// GetAllAportesByID returns all aportes in the fondo from a specific user
func (h *UsersHandler) GetAllAportesByID(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	id := getID(r)
	h.l.Info("[GetAllAportesByID] Recieving call to get all aportes from", "user", us)

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	aportes, err := h.UserService.GetAllAportesByID(id, startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&aportes, w)
}

// GetSumAportesByID returns sum of aportes of an specific user
func (h *UsersHandler) GetSumAportesByID(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	id := getID(r)

	h.l.Info("[GetSumAportesByID] Recieving call to get sum of aportes from", "user", us)
	aportes, err := h.UserService.GetSumAportesByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&aportes, w)
}

// GetProyeccionCredito returns an arary of cuotas based on the given credit
func (h *UsersHandler) GetProyeccionCredito(w http.ResponseWriter, r *http.Request) {
	var us = (context.Get(r, "us")).(data.User)
	var cr = (context.Get(r, "cr")).(*data.Credito)
	h.l.Info("[CalculateCredito] Recieving call to get cuotas from ", "user", us)
	cuotas := h.UserService.CalcularCredito(cr)
	data.ToJSON(&cuotas, w)
}
