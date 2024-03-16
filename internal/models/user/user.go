package user_model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
	Role     string
}

type UserLogin struct {
	Login    string `json:"login" valid:",required"`
	Password string `json:"password" valid:",required"`
}
