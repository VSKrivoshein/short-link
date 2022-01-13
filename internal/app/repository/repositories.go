package repositories

import (
	author_repo "github.com/VSKrivoshein/short-link/internal/app/repository/author"
	shortener_repo "github.com/VSKrivoshein/short-link/internal/app/repository/shortener"
	"github.com/VSKrivoshein/short-link/internal/app/services/author"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repositories struct {
	Shortener shortener.Repository
	Author    author.Repository
}

func NewRepositories(dbURL string) *Repositories {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logrus.Fatalf("Postgres connection failed %s", err.Error())
	}

	logrus.Info("DB is connected")

	return &Repositories{
		Shortener: shortener_repo.NewRepository(db),
		Author:    author_repo.NewRepository(db),
	}
}
