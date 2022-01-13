package services

import (
	repositories "github.com/VSKrivoshein/short-link/internal/app/repository"
	"github.com/VSKrivoshein/short-link/internal/app/services/author"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
)

type Services struct {
	Shortener shortener.Service
	Author    author.Service
}

func NewServices(r *repositories.Repositories) *Services {
	return &Services{
		Shortener: shortener.NewService(r.Shortener),
		Author:    author.NewService(r.Author),
	}
}
