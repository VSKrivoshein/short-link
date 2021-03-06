package author

import (
	"net/http"
	"time"

	"github.com/VSKrivoshein/short-link/internal/app/e"
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 14
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", e.New(err, e.ErrToken, http.StatusInternalServerError)
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetExpirationTime() time.Time {
	return time.Now().Add(time.Hour * 3)
}
