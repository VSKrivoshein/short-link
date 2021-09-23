package main

import (
	cnf "github.com/VSKrivoshein/short-link/internal/app/configs"
	u "github.com/VSKrivoshein/short-link/internal/app/utils"
	"github.com/sirupsen/logrus"
)

// @title Short link with authorization
// @version 1.0
// @description Service of the short link with authorization, hexagonal architecture, integration test
// @host localhost:8080
func main() {
	u.UseJSONLogFormat()

	handler := cnf.InitServices()

	errors := make(chan error, 2)
	go func() {
		errors <- cnf.StartServices(handler)
	}()

	go func() {
		errors <- cnf.GracefulShutdown()
	}()

	logrus.Fatalf("Terminated %s", <-errors)
}
