package handlers

import (
	"authentication-api/data"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Signup hanldes user signup requests
func (h *Auth) Signup(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[Signup] Handling signup request")

	user := r.Context().Value(KeyUser{}).(*data.UserCreate)
	err := h.u.CreateUser(user)

	if err != nil {
		h.l.Error("[Signup] Something went wrong creating an user in the database ", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}

	data.ToJSON(&user, w)
}

// Signin hanldes user Signin requests
func (h *Auth) Signin(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(KeyUser{}).(*data.UserSignin)

	h.l.Info("[Signin] Handling signin request for", "user", user.Usuario)

	userdb, err := h.u.GetUserByUser(user.Usuario)

	switch err {
	case nil:
		err = bcrypt.CompareHashAndPassword([]byte(userdb.Contrasena), []byte(user.Contrasena))
		if err != nil {
			h.l.Info("[Signin] Failed login attempt at", "email", user.Email)
			data.ToJSON(&GenericError{Message: "Can't authenticate user"}, w)
			return
		}
		tokenString, err := h.GenerateToken(&userdb)
		if err != nil {
			h.l.Error("[Signin] Something went wrong generating token", "error", err)
			data.ToJSON(&GenericError{Message: "Something went wrong generating token"}, w)
			return
		}
		data.ToJSON(&Token{Message: tokenString}, w)
	case data.ErrProductNotFound:
		w.WriteHeader(http.StatusNotFound)
		h.l.Error("[Signin] User with phone not found", "email", user.Email)
		data.ToJSON(&GenericError{Message: err.Error()}, w)

		return
	default:
		h.l.Error("[Signin] Fetching user with email", "email", user.Email, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
