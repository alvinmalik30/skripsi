package security

import "github.com/golang-jwt/jwt/v5"

type AppClaims struct {
	jwt.RegisteredClaims
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Role     string   `json:"role,omitempty"`
	Services []string `json:"services,omitempty"`
}

// Service -> Resource/API apa yg user boleh access (Authorization)
