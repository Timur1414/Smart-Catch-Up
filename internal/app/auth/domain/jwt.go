package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	Uuid      string
	UserId    int
	ExpiredAt time.Time
}

type Claims struct {
	jwt.RegisteredClaims
	Type string `json:"typ,omitempty"`
}
