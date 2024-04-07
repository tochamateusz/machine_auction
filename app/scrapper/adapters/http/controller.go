package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Register(r *gin.Engine) {
	scrappe := r.Group("scrapper")
	scrappe.POST("start", Scrap)

}

func Scrap(ctx *gin.Context) {
	log.Info().Msg("Scrapping started...")
}
