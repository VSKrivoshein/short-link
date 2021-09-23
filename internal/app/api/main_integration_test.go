package api_test

import (
	cnf "github.com/VSKrivoshein/short-link/internal/app/configs"
	"net/http/httptest"
	"os"
	"testing"
)

var TestSrv *httptest.Server

func TestMain(m *testing.M) {
	srv := cnf.InitServices().InitRoutes(os.Getenv("SRV_GIN_MODE"))
	TestSrv = httptest.NewServer(srv)
	defer TestSrv.Close()

	exitVal := m.Run()

	os.Exit(exitVal)
}
