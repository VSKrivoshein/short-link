package auther

import (
	"errors"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var jwtKey = []byte("the_most_secret_key")

type service struct {
	repo      Repository
	validator *validator.Validate
}

func NewService(repo Repository) Service {
	return &service{
		repo:      repo,
		validator: validator.New(),
	}
}

type Claims struct {
	UserId string
	jwt.StandardClaims
}

// todo make validation error correct

func (s *service) SingIn(user *User) error {
	if err := s.validator.Struct(user); err != nil {
		return e.New(err, err, http.StatusUnprocessableEntity)
	}

	if err := s.repo.GetUser(user); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	if !CheckPasswordHash(user.Password, user.PasswordHash) {
		return e.New(
			errors.New("if !CheckPasswordHash(user.Password, user.PasswordHash)"),
			e.ErrCredentialsInvalid,
			http.StatusUnauthorized,
		)
	}

	user.TokenExpiration = GetExpirationTime()

	claims := &Claims{
		UserId: user.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: user.TokenExpiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return e.New(err, e.ErrToken, http.StatusInternalServerError)
	}

	user.TokenString = tokenString

	return nil
}

func (s *service) SingUp(user *User) error {
	if err := s.validator.Struct(user); err != nil {
		return e.New(err, err, http.StatusUnprocessableEntity)
	}

	passwordHash, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	user.PasswordHash = passwordHash

	if err := s.repo.CreateUser(user); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	user.TokenExpiration = GetExpirationTime()

	claims := &Claims{
		UserId: user.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: user.TokenExpiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return e.New(err, e.ErrToken, http.StatusInternalServerError)
	}

	user.TokenString = tokenString

	return nil
}

func (s *service) DeleteUser(user *User) error {
	if err := s.repo.DeleteUser(user); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}
	return nil
}

func (s *service) CheckAuthAndRefresh(user *User) error {

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(
		user.TokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		return e.New(err, e.ErrToken, http.StatusInternalServerError)
	}

	if !tkn.Valid {
		return e.New(errors.New("if !tkn.Valid"), e.ErrToken, http.StatusForbidden)
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 2*time.Hour {
		return nil
	}

	user.TokenExpiration = GetExpirationTime()
	claims.ExpiresAt = user.TokenExpiration.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return e.New(err, e.ErrToken, http.StatusInternalServerError)
	}

	user.TokenString = tokenString
	user.UserId = claims.UserId

	return nil
}
