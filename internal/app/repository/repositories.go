package repositories

import (
	auther_repo "github.com/VSKrivoshein/short-link/internal/app/repository/auther"
	shortener_repo "github.com/VSKrivoshein/short-link/internal/app/repository/shortener"
	"github.com/VSKrivoshein/short-link/internal/app/services/auther"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repositories struct {
	Shortener shortener.Repository
	Auther    auther.Repository
}

func NewRepositories(dbUrl string) *Repositories {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		logrus.Fatalf("Postgres connection failed %s", err.Error())
	}

	logrus.Info("DB is connected")

	return &Repositories{
		Shortener: shortener_repo.NewRepository(db),
		Auther:    auther_repo.NewRepository(db),
	}
}
