package services

import (
	repositories "github.com/VSKrivoshein/short-link/internal/app/repository"
	"github.com/VSKrivoshein/short-link/internal/app/services/auther"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
)

type Services struct {
	Shortener shortener.Service
	Auther     auther.Service
}

func NewServices(repositories *repositories.Repositories) *Services {
	return &Services{
		Shortener: shortener.NewService(repositories.Shortener),
		Auther:     auther.NewService(repositories.Auther),
	}
}
