package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	id            int
	name          string
	password_hash string
}

type Entrance struct {
	Email    string `json:"email"`
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

type RecoveryPasswordData struct { // структура для получения контакных данных для отправки кода
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type DataRecoveryPassword struct { // структура для получения контакных данных и сверкой кода
	Email string `json:"email"`
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type RecoveryPassword struct { // структура для восстановления пароля
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	Token                string `json:"token"`
	ConfirmationPassword string `json:"confirmation_password"`
	NewPassword          string `json:"new_password"`
}
