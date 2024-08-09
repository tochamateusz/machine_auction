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

type AuctionEvents struct{ name string }

func (a AuctionEvents) Name() string {
	return a.name
}

var (
	Founded = AuctionEvents{name: "auction.founded"}
)

type HttpScrapperApi struct {
	scrapper   *auctionScrapper.Scrapper
	eventBus   events.IEventBus
	repository auction.Repository
}

func Init(r *gin.Engine) {
	scrapperGroup := r.Group("scrapper")

	repository, err := auction_file.NewFileAuctionRepository()
	if err != nil {
		log.Fatal().Err(err).Msgf("%p", err)
	}

	scrapper, err := auctionScrapper.NewScrapper(repository)
	if err != nil {
		log.Debug().Caller().Msgf("%p", err)
		panic(err)
	}

	eventBus := events.NewEventBus()
	eventBus.Listen(Founded.Name(), scrapper.OnAuctionFound)
	// eventBus.Listen("auctions.founded", func(ctx context.Context, message interface{}) {

	// 	scrappingRepository, err := scrapping_events.NewFileScrappedAuctionsRepository()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	auction, ok := message.(acutions_events.AuctionsFounded)
	// 	if ok == false {
	// 		log.Warn().Msgf("Malformed message %+v", message)
	// 		return
	// 	}
	// 	for _, v := range auction.Auction {
	// 		v.CreatedAtDate(time.Now().Format(time.DateTime))
	// 		scrappingRepository.Save(scrapping.ScrappedAuctions{
	// 			Id:        auction.Id,
	// 			Triggered: time.Now().Format(time.RFC3339Nano),
	// 			Auction:   v,
	// 		})
	// 	}
	// })

	go eventBus.Serve(context.Background())

	http_client := HttpScrapperApi{
		scrapper:   scrapper,
		repository: repository,
		eventBus:   eventBus,
	}

	scrapperGroup.POST("start", http_client.BaseScrap)
	scrapperGroup.GET(":id", http_client.Get)
	scrapperGroup.GET("/images/:id", http_client.GetAllImage)
	scrapperGroup.GET("", http_client.GetAll)
	scrapperGroup.GET("test", http_client.TEST)
}

type HttpError struct {
	Reason string `json:"reason"`
}

func (h *HttpScrapperApi) BaseScrap(ctx *gin.Context) {
	auctions, err := h.scrapper.GetAuctions()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, HttpError{})
		return
	}

	mapAuctions := h.scrapper.Scrap(ctx, auctions)

	for _, auction := range mapAuctions {
		h.eventBus.Dispatch(Founded.Name(), auction)
	}
	ctx.JSON(http.StatusOK, mapAuctions)
}

func (h *HttpScrapperApi) GetAllImage(ctx *gin.Context) {
	id := ctx.Param("id")
	auction, err := h.scrapper.GetAllImages(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	log.Info().Msgf("GET: auction: %+v\n", auction)
	ctx.JSON(http.StatusOK, auction)
}

func (h *HttpScrapperApi) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	auction := h.repository.Get(id)
	log.Info().Msgf("GET: auction: %+v\n", auction)
	ctx.JSON(http.StatusOK, AuctionDTO{
		Id:            auction.Id(),
		Image:         auction.Image(),
		Name:          auction.Name(),
		Year:          auction.Year(),
		Price:         auction.Price(),
		EndDate:       auction.EndDate(),
		Description:   auction.Description(),
		StartingPrice: auction.StartingPrice(),
	})
}

type AuctionDTO struct {
	Id            string   `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Year          string   `json:"year"`
	Price         string   `json:"price"`
	EndDate       string   `json:"end_date"`
	Description   []string `json:"description"`
	StartingPrice string   `json:"starting_price"`
}

func (h *HttpScrapperApi) TEST(ctx *gin.Context) {

	for _, v := range []string{"test", "test1", "test2", "test3"} {
		auction := h.repository.Get("10000")
		if auction.Id() != "10000" {
			return
		}
		auction.Describe([]string{v})
		go h.repository.Save(auction)
	}

}

func (h *HttpScrapperApi) GetAll(ctx *gin.Context) {
	auctions := h.repository.GetAll()
	var auctionsDtos []AuctionDTO
	for _, v := range auctions {
		auctionDto := AuctionDTO{
			Id:          v.Id(),
			Image:       v.Image(),
			Name:        v.Name(),
			Year:        v.Year(),
			Price:       v.Price(),
			EndDate:     v.EndDate(),
			Description: v.Description(),
		}
		auctionsDtos = append(auctionsDtos, auctionDto)
	}
	ctx.JSON(http.StatusOK, auctionsDtos)
}
