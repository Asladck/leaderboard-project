package domain

import "github.com/google/uuid"

// User represents an application user.
type User struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
}
