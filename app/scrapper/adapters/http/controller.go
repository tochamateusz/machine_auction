package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	auctionScrapper "github.com/tochamateusz/machine_auction/app/scrapper"
	"github.com/tochamateusz/machine_auction/domain/auction"
	"github.com/tochamateusz/machine_auction/infrastructure/events"
	auction_file "github.com/tochamateusz/machine_auction/infrastructure/repository/auction"
)

type Key = string

type HttpScrapperApi struct {
	scrapper   *auctionScrapper.Scrapper
	repository auction.Repository
	eventBus   events.IEventBus
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

	eventBus := events.NewEventBus()
	eventBus.Listen("test", func(ctx context.Context, message interface{}) {
		log.Info().Msg("ON test message")
	})

	go eventBus.Serve(context.Background())

	repository, err := auction_file.NewFileAuctionRepository()
	if err != nil {
		log.Fatal().Msgf("%p", err)
	}

	for _, v := range auctions {
		repository.Save(v)
	}

	http_client := HttpScrapperApi{
		scrapper:   scrapper,
		repository: repository,
		eventBus:   eventBus,
	}

	scrapper.SaveImage("")

	scrapperGroup.POST("start", http_client.BaseScrap)
	scrapperGroup.GET(":id", http_client.Get)
	scrapperGroup.GET("", http_client.GetAll)

}

func (h *HttpScrapperApi) BaseScrap(ctx *gin.Context) {
	auctions, err := h.scrapper.GetAuctions()
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	for _, v := range auctions {
		h.repository.Save(v)
		h.eventBus.Dispatch("test", struct{ Test string }{Test: "test"})

	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func (h *HttpScrapperApi) Get(ctx *gin.Context) {
	auction := h.repository.Get("9768")
	log.Info().Msgf("GET: auction: %+v\n", auction)
	ctx.JSON(http.StatusOK, auction)
}

type AuctionDTO struct {
	Id      string `json:"id"`
	Image   string `json:"image"`
	Name    string `json:"name"`
	Year    string `json:"year"`
	Price   string `json:"price"`
	EndDate string `json:"end_date"`
}

func (h *HttpScrapperApi) GetAll(ctx *gin.Context) {
	auctions := h.repository.GetAll()
	var auctionsDtos []AuctionDTO
	for _, v := range auctions {
		auctionDto := AuctionDTO{
			Id:      v.Id(),
			Image:   v.Image(),
			Name:    v.Name(),
			Year:    v.Year(),
			Price:   v.Price(),
			EndDate: v.EndDate(),
		}
		auctionsDtos = append(auctionsDtos, auctionDto)
	}
	ctx.JSON(http.StatusOK, auctionsDtos)
}
