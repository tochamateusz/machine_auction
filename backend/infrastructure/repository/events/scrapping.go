package events

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/domain/auction"
	"github.com/tochamateusz/machine_auction/domain/scrapping"
	auction_file "github.com/tochamateusz/machine_auction/infrastructure/repository/auction"
)

const storePath = "./db/scrapping-events.json"

type FileScrappedAuctionsRepository struct {
	file *os.File
}

// GetAll implements scrapping.Repository.
func (*FileScrappedAuctionsRepository) GetAll() []scrapping.ScrappedAuctions {
	panic("unimplemented")
}

// Get implements scrapping.Repository.
func (*FileScrappedAuctionsRepository) Get(id string) []scrapping.ScrappedAuctions {
	raw, err := os.ReadFile(storePath)
	if err != nil {
		log.Info().Msgf("Can't read a file %s", err.Error())
		return []scrapping.ScrappedAuctions{}
	}
	scrappedAuctionsDataModel := []ScrappedAuctionsDataModel{}

	json.Unmarshal(raw, &scrappedAuctionsDataModel)

	founded := []scrapping.ScrappedAuctions{}
	for _, v := range scrappedAuctionsDataModel {
		if v.Id == id {
			founded = append(founded, scrapping.ScrappedAuctions{
				Id:        v.Id,
				Triggered: v.Triggered,
				Auction: *auction.NewAuction(
					v.Auction.Id,
					v.Auction.Id,
					v.Auction.Name,
					v.Auction.Year,
					v.Auction.Price,
					v.Auction.EndDate),
			})
		}
	}

	return founded
}

type ScrappedAuctionsDataModel struct {
	Id        string
	Triggered string
	Auction   auction_file.AuctionDataModel
}

// Save implements scrapping.Repository.
func (f *FileScrappedAuctionsRepository) Save(auction scrapping.ScrappedAuctions) {
	m := sync.RWMutex{}
	m.Lock()
	defer m.Unlock()

	raw, _ := os.ReadFile(storePath)
	scrappedAuctionsDataModel := []ScrappedAuctionsDataModel{}

	json.Unmarshal(raw, &scrappedAuctionsDataModel)

	auctionDataModel := auction_file.AuctionDataModel{
		Id:      auction.Auction.Id(),
		Image:   auction.Auction.Image(),
		Name:    auction.Auction.Name(),
		Year:    auction.Auction.Year(),
		Price:   auction.Auction.Price(),
		EndDate: auction.Auction.EndDate(),
	}

	scrappedAuctionsDataModel = append(scrappedAuctionsDataModel, ScrappedAuctionsDataModel{
		Id:        auction.Id,
		Triggered: auction.Triggered,
		Auction:   auctionDataModel,
	})

	bytes, err := json.Marshal(scrappedAuctionsDataModel)
	if err != nil {
		return
	}

	f.file.Seek(0, 0)
	_, err = f.file.Write(bytes)
	if err != nil {
		log.Err(err).Msg("cannnot write auction store")
		return
	}

}

func NewFileScrappedAuctionsRepository() (scrapping.Repository, error) {
	file, err := os.OpenFile(storePath, os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		return nil, fmt.Errorf("cannnot create scrapper")
	}
	return &FileScrappedAuctionsRepository{
		file,
	}, nil
}
