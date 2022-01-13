package shortener

import (
	"fmt"
	"net/http"

	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/go-playground/validator"
	"github.com/teris-io/shortid"
)

type service struct {
	redirectRepo Repository
	validator    *validator.Validate
}

func NewService(redirectRepo Repository) Service {
	return &service{
		redirectRepo,
		validator.New(),
	}
}

func (s *service) GetLink(redirect *Redirect) error {
	err := s.redirectRepo.GetLink(redirect)
	if err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}
	return nil
}

func (s *service) GetAllLinks(redirect *Redirect) error {
	err := s.redirectRepo.GetAllLinks(redirect)
	if err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}
	return nil
}

func (s *service) CreateLink(redirect *Redirect) error {
	if err := s.validator.Struct(redirect); err != nil {
		return e.New(err, err, http.StatusUnprocessableEntity)
	}

	redirect.LinkHash = shortid.MustGenerate()

	if err := s.redirectRepo.CreateLink(redirect); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	return nil
}

func (s *service) DeleteLink(redirect *Redirect) error {
	err := s.redirectRepo.DeleteLink(redirect)
	if err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}
	return nil
}
