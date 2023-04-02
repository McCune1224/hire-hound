package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`

	Applications []Application `json:"applications"`
}

type UserContextJWT struct {
	ID       uint `json:"id"`
	Username string  `json:"username"`
}
