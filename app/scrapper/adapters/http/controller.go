package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auctionScrapper "github.com/tochamateusz/machine_auction/app/scrapper"
)

type Key = string

type HttpScrapperApi struct {
	scrapper auctionScrapper.Scrapper
}

func Register(r *gin.Engine) {
	scrapperGroup := r.Group("scrapper")
	scrapper := auctionScrapper.Scrapper{}

	scrapper.Login()
	scrapper.PrintCookie()

	http_client := HttpScrapperApi{
		scrapper,
	}

	scrapperGroup.POST("start", http_client.Login)

}

func (h *HttpScrapperApi) Login(ctx *gin.Context) {
	err := h.scrapper.Login()
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}
