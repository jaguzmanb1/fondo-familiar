package data

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// User define la estructura de un usuario para el API
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

// UserSignin defines user when is on Signin phase
type UserSignin struct {
	ID         int    `json:"id"`
	Nombre     string `json:"nombre"`
	Usuario    string `json:"usuario" validate:"required"`
	Email      string `json:"email"`
	Contrasena string `json:"contrasena" validate:"required"`
	IDRol      int    `json:"idRol"`
}

// UserCreate defines data user structure when realices a signup
type UserCreate struct {
	ID         int    `json:"id"`
	Nombre     string `json:"nombre" validate:"required"`
	Celular    string `json:"celular"`
	Contrasena string `json:"contrasena" validate:"required"`
	Email      string `json:"email"`
	Usuario    string `json:"usuario" validate:"required"`
}

//Users is una colección de User
type Users []*User

//UserService representa una implementación de mysql
type UserService struct {
	DB *sql.DB
	l  hclog.Logger
}

// New creates a new user service
func New(d *sql.DB, l hclog.Logger) *UserService {
	return &UserService{d, l}
}

// GetUsers retorna una lista de usuarios
func (s *UserService) GetUsers() (Users, error) {
	users := Users{}
	rows, err := s.DB.Query("SELECT id, nombre FROM usuario")
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

//GetUserByID returns an user given an id
func (s *UserService) GetUserByID(id int) (User, error) {
	user := User{}
	rows, err := s.DB.Query("SELECT id, nombre FROM usuario WHERE id = (?)", id)
	if err != nil {
		return user, ErrProductNotFound
	}

	for rows.Next() {
		user = User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return user, err
		}

		return user, err
	}

	return user, ErrProductNotFound
}

//GetUserByEmail returns an user given an email
func (s *UserService) GetUserByEmail(email string) (UserSignin, error) {
	s.l.Info("[GetUserByEmail] Getting user from database with", "email", email)

	user := UserSignin{}
	rows, err := s.DB.Query("SELECT id, nombre, contrasena, idRol, email FROM usuario WHERE email = (?)", email)
	if err != nil {
		return user, ErrProductNotFound
	}

	for rows.Next() {
		user = UserSignin{}
		err = rows.Scan(&user.ID, &user.Nombre, &user.Contrasena, &user.IDRol, &user.Email)
		if err != nil {
			return user, err
		}

		return user, err
	}

	return user, ErrProductNotFound
}

//GetUserByUser returns an user given an user
func (s *UserService) GetUserByUser(usuario string) (UserSignin, error) {
	s.l.Info("[GetUserByEmail] Getting user from database with", "user", usuario)

	user := UserSignin{}
	rows, err := s.DB.Query("SELECT id, nombre, contrasena, idRol, email, usuario FROM usuario WHERE usuario = (?)", usuario)
	if err != nil {
		return user, ErrProductNotFound
	}

	for rows.Next() {
		user = UserSignin{}
		err = rows.Scan(&user.ID, &user.Nombre, &user.Contrasena, &user.IDRol, &user.Email, &user.Usuario)
		if err != nil {
			return user, err
		}

		return user, err
	}

	return user, ErrProductNotFound
}

//CreateUser crea un usuario
func (s *UserService) CreateUser(pUser *UserCreate) error {
	s.l.Info("[CreateUser] Creating", "user", pUser)
	saltedPassword, err := s.hashAndSalt([]byte(pUser.Contrasena))
	if err != nil {
		return err
	}
	_, err = s.DB.Exec("INSERT INTO usuario (nombre, celular, contrasena, email, idrol, usuario) VALUES (?, ?, ?, ?, ?, ?)",
		pUser.Nombre,
		pUser.Celular,
		saltedPassword,
		pUser.Email,
		3,
		pUser.Usuario)

	return err
}

//DeleteUser elimina un usuario dado un id
func (s *UserService) DeleteUser(id int) error {
	_, err := s.DB.Exec("DELETE FROM usuario WHERE id = (?)", id)
	if err != nil {
		return err
	}
	return nil
}
