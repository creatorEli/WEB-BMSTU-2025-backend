package ds

import (
	"time_of_armies/internal/app/role"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID uuid.UUID `json:"user_uuid"`
	Scopes   []string  `json:"scopes" json:"scopes"` // список доступов в нашей системе
	Login    string
	Role     role.Role
}
