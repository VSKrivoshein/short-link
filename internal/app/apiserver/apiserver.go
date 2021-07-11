package apiserver

import (
	"github.com/VSKrivoshein/short-link/internal/app/store"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *gin.Engine
	store  *store.Store
}

func NewAPIServer(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: gin.New(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

// todo refactor
func (s *APIServer) configureRouter() {
	s.router.GET("/hello", s.handleHello())
}

func (s *APIServer) configureStore() error {
	st := store.NewStore(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) handleHello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	}
}

