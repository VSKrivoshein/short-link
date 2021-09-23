package api

import (
	"fmt"
	_ "github.com/VSKrivoshein/short-link/docs"
	"github.com/VSKrivoshein/short-link/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
)

type Config struct {
	Protocol string
	Host     string
	Port     string
	GinMode  string
	LogLvl   string
	Origin   string
}

type Handler struct {
	Services *services.Services
	Config   Config
}

func NewHandler(services *services.Services) *Handler {
	config := Config{
		Protocol: os.Getenv("SRV_PROTOCOL"),
		Host:     os.Getenv("SRV_HOST"),
		Port:     os.Getenv("SRV_PORT"),
		GinMode:  os.Getenv("SRV_GIN_MODE"),
		LogLvl:   os.Getenv("SRV_LOG_LVL"),
		Origin: fmt.Sprintf(
			"%v://%v:%v/",
			os.Getenv("SRV_PROTOCOL"),
			os.Getenv("SRV_HOST"),
			os.Getenv("SRV_PORT")),
	}

	return &Handler{services, config}
}

func (h *Handler) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(JSONLogMiddleware())
	r.Use(ErrorMiddleware())

	r.GET("/health", h.health)
	r.GET("/:hash", h.redirect)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-out", h.signOut)
		auth.DELETE("/delete-user", h.deleteUser)
	}

	links := r.Group("/links", h.checkAuthAndRefreshMiddleware)
	{
		links.GET("/get-all", h.getAllLinks)
		links.POST("/create", h.create)
		links.DELETE("/delete", h.deleteLink)
	}

	return r
}
