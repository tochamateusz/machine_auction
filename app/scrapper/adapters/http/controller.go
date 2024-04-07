package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	auctionScrapper "github.com/tochamateusz/machine_auction/app/scrapper"
	"github.com/tochamateusz/machine_auction/domain/auction"
	auction_file "github.com/tochamateusz/machine_auction/infrastructure/repository/auction"
)

type Key = string

type HttpScrapperApi struct {
	scrapper   *auctionScrapper.Scrapper
	repository auction.Repository
}

func Init(r *gin.Engine) {
	scrapperGroup := r.Group("scrapper")
	scrapper, err := auctionScrapper.NewScrapper()
	if err != nil {
		log.Fatal().Msgf("%p", err)
	}

	scrapper.Login()
	scrapper.PrintCookie()
	auctions, err := scrapper.GetAuctions()
	if err != nil {
		log.Fatal().Msgf("%p", err)
	}

	repository, err := auction_file.NewFileAuctionRepository()
	if err != nil {
		log.Fatal().Msgf("%p", err)
	}

	for _, v := range auctions {
		repository.Save(v)
	}

	http_client := HttpScrapperApi{
		scrapper,
		repository,
	}

	scrapperGroup.POST("start", http_client.Login)
	scrapperGroup.GET(":id", http_client.Get)
	scrapperGroup.GET("", http_client.GetAll)

}

func (h *HttpScrapperApi) Login(ctx *gin.Context) {
	err := h.scrapper.Login()
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func (h *HttpScrapperApi) Get(ctx *gin.Context) {
	auction := h.repository.Get("9768")
	log.Info().Msgf("GET: auction: %+v\n", auction)
	ctx.JSON(http.StatusOK, auction)
}

func (h *HttpScrapperApi) GetAll(ctx *gin.Context) {
	auctions := h.repository.GetAll()
	log.Info().Msgf("GET: auction: %+v\n", auctions)
	ctx.JSON(http.StatusOK, auctions)
}
