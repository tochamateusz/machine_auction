package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) DbHealthCheck(c *gin.Context) {
	log.Info().Msg("Pong")
	c.Status(http.StatusOK)
}
