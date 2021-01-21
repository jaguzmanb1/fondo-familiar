package auth

import (
	"fmt"
	"fondo-mod/data"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
)

var signKey = []byte(os.Getenv("secretAuthKey"))

// Auth define a object of auth to validate requests tokens
type Auth struct {
	l hclog.Logger
}

// AuthError is a generic auth error message returned by a server
type AuthError struct {
	Message string `json:"message"`
}

// New creates a new auth validator instance
func New(l hclog.Logger) *Auth {
	l.Debug("[New] Creating new auth instance")

	return &Auth{l}
}

// KeyClient usada para el middleware
type KeyClient struct{}

// ValidateToken validates a user request to be signed with a JWT token
func (h *Auth) validateToken(t string, l string) (bool, data.User, error) {
	h.l.Info("[validateToken] Validating token")

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[validateToken] Unexpected signing method: %v", token.Header["alg"])
		}

		return signKey, nil
	})

	if err != nil {
		return false, data.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		rol := claims["rol"].(float64)
		email := claims["email"].(string)
		id := claims["id"].(float64)

		s := fmt.Sprintf("%.0f", rol)
		if s <= l {
			return true, data.User{ID: int(id), Rol: s, Email: email}, nil
		}

	}

	return false, data.User{}, err

}
