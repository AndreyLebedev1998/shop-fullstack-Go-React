package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	id            int
	name          string
	password_hash string
}

type Entrance struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Credentials struct {
	Name     string  `json:"name"`
	LastName string  `json:"lastname"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

type UserEntrance struct {
	Name     string  `json:"name"`
	LastName string  `json:"lastname"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
