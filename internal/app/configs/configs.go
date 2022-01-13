package cnf

import (
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	repositories "github.com/VSKrivoshein/short-link/internal/app/repository"
	"github.com/VSKrivoshein/short-link/internal/app/services"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func InitServices() *api.Handler {
	var dbURL = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)

	repos := repositories.NewRepositories(dbURL)
	service := services.NewServices(repos)

	return api.NewHandler(service)
}

func StartServices(handler *api.Handler) error {
	port := fmt.Sprintf(":%s", handler.Config.Port)
	logrus.Infof("Starting server on port %s", port)
	return http.ListenAndServe(port, handler.InitRoutes(handler.Config.GinMode))
}

func GracefulShutdown() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	return fmt.Errorf("%s", <-signals)
}
