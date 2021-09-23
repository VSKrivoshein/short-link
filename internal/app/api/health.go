package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Health check
// @Tags health
// @Description Check that service alive and ready for traffic
// @Success 200 {string} string "I am alive"
// @Failure 500
// @Router /health [get]
func (h *Handler) health(c *gin.Context) {
	c.String(http.StatusOK, "I am alive")
}
