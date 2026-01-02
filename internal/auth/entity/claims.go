package entity

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   string         `json:"user_id"`
	UserType string         `json:"user_type"`
	Custom   map[string]any `json:"custom"`
}
