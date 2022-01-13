package author

import (
	"time"
)

type User struct {
	UserID          string `db:"id"`
	Email           string `validate:"required,email" db:"email"`
	Password        string `validate:"min=6,max=100"`
	PasswordHash    string `db:"password_hash"`
	TokenString     string
	TokenExpiration time.Time
}
