package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	auctionScrapper "github.com/tochamateusz/machine_auction/app/scrapper"
	"github.com/tochamateusz/machine_auction/domain/auction"
	"github.com/tochamateusz/machine_auction/domain/scrapping"
	"github.com/tochamateusz/machine_auction/infrastructure/events"
	auction_file "github.com/tochamateusz/machine_auction/infrastructure/repository/auction"
	scrapping_events "github.com/tochamateusz/machine_auction/infrastructure/repository/events"

	acutions_events "github.com/tochamateusz/machine_auction/domain/auction"
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
		log.Debug().Caller().Msgf("%p", err)
		panic(err)
	}

	scrapper.Login()
	scrapper.PrintCookie()

	scrappingRepository, err := scrapping_events.NewFileScrappedAuctionsRepository()
	if err != nil {
		panic(err)
	}

	eventBus := events.NewEventBus()
	eventBus.Listen("auction.founded", scrapper.OnAuctionFound)
	eventBus.Listen("auctions.founded", func(ctx context.Context, message interface{}) {
		auction, ok := message.(acutions_events.AuctionsFounded)
		if ok == false {
			log.Warn().Msgf("Malformed message %+v", message)
			return
		}
		for _, v := range auction.Auction {
			scrappingRepository.Save(scrapping.ScrappedAuctions{
				Id:        auction.Id,
				Triggered: time.Now().Format(time.RFC3339Nano),
				Auction:   v,
			})
		}
	})

	go eventBus.Serve(context.Background())

	repository, err := auction_file.NewFileAuctionRepository()
	if err != nil {
		log.Fatal().Msgf("%p", err)
	}

	http_client := HttpScrapperApi{
		scrapper:   scrapper,
		repository: repository,
		eventBus:   eventBus,
	}

	scrapperGroup.POST("start", http_client.BaseScrap)
	scrapperGroup.GET(":id", http_client.Get)
	scrapperGroup.GET("/images/:id", http_client.GetAllImage)
	scrapperGroup.GET("", http_client.GetAll)

}

func (h *HttpScrapperApi) BaseScrap(ctx *gin.Context) {
	auctions, err := h.scrapper.GetAuctions()
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	mapAuctions := make(map[string]bool)
	mutex := sync.Mutex{}

	h.scrapper.RegisterDone(func(done string) {
		mutex.Lock()
		defer mutex.Unlock()
		mapAuctions[done] = true
		if len(mapAuctions) == len(auctions) {
			log.Info().Msgf("All scrapped jobs done")
		}
	})

	h.eventBus.Dispatch("auctions.founded", auction.AuctionsFounded{
		Id:      uuid.NewString(),
		Auction: auctions,
	})

	for _, v := range auctions {
		h.repository.Save(v)
		h.eventBus.Dispatch("auction.founded", auction.AuctionFounded{
			Auction: v,
		})
	}

	ctx.Writer.WriteHeader(http.StatusOK)
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
		Id:      auction.Id(),
		Image:   auction.Image(),
		Name:    auction.Name(),
		Year:    auction.Year(),
		Price:   auction.Price(),
		EndDate: auction.EndDate(),
	})
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