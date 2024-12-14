package models

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	UserRole string `json:"user_role"`
}