package dto

import "github.com/google/uuid"

type UserDTO struct {
	UUID     uuid.UUID `json:"uuid,omitempty"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Position string    `json:"position"`
	Team     string    `json:"team"`
	Role     string    `json:"role,omitempty"`
	Password string    `json:"password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
