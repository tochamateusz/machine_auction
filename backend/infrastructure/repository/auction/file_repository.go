package auction_file

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

const storePath = "./db/auctions.json"

type FileAuctionRepository struct {
	file    *os.File
	mu      *sync.Mutex
	version uint
}

// Get implements AuctionRepository.
func (f *FileAuctionRepository) Get(id string) auction.Auction {
	f.mu.Lock()
	defer f.mu.Unlock()
	raw, err := os.ReadFile(storePath)
	if err != nil {
		log.Info().Msgf("Can't read a file %s", err.Error())
		return auction.Auction{}
	}

	auctionDataModel := make(map[string]AuctionDataModel)
	json.Unmarshal(raw, &auctionDataModel)

	v, ok := auctionDataModel[id]
	if ok == false {
		log.Info().Msgf("auction id=%v not exist", ok)
		return auction.Auction{}
	}
	a := *auction.NewAuction(v.Id, v.Image, v.Name, v.Year, v.Price, v.EndDate)
	a.Describe(v.Description)
	return a
}

// GetAll implements AuctionRepository.
func (f *FileAuctionRepository) GetAll() []auction.Auction {
	f.mu.Lock()
	defer f.mu.Unlock()
	raw, err := os.ReadFile(storePath)
	if err != nil {
		log.Err(err).Msgf("GetAll")
		return []auction.Auction{}
	}
	auctionDataModel := make(map[string]AuctionDataModel)

	json.Unmarshal(raw, &auctionDataModel)
	var array []auction.Auction
	for _, v := range auctionDataModel {
		a := *auction.NewAuction(v.Id, v.Image, v.Name, v.Year, v.Price, v.EndDate)
		a.Describe(v.Description)
		array = append(array, a)
	}
	return array
}

type AuctionDataModel struct {
	Id            string   `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Year          string   `json:"year"`
	Price         string   `json:"price"`
	EndDate       string   `json:"end_date"`
	CreatedAt     string   `json:"create_at"`
	Description   []string `json:"description,omitempty"`
	StartingPrice string   `json:"starting_price"`
}

// Save implements AuctionRepository.
func (f *FileAuctionRepository) Save(a auction.Auction) {
	f.mu.Lock()
	defer func() {
		f.mu.Unlock()
	}()

	f.version++

	raw, err := os.ReadFile(storePath)
	if err != nil {
		log.Err(err).Msgf("save")
		return
	}
	auctionDataModel := make(map[string]AuctionDataModel)

	err = json.Unmarshal(raw, &auctionDataModel)
	if err != nil {
		log.Err(err).Msgf("save")
		return
	}

	auctionDataModel[a.Id()] = AuctionDataModel{
		Id:            a.Id(),
		Image:         a.Image(),
		Name:          a.Name(),
		Year:          a.Year(),
		Price:         a.Price(),
		EndDate:       a.EndDate(),
		Description:   a.Description(),
		StartingPrice: a.StartingPrice(),
	}

	bytes, err := json.Marshal(auctionDataModel)
	if err != nil {
		return
	}

	f.file.Truncate(0)
	f.file.Seek(0, 0)
	if err != nil {
		log.Err(err).Msg("cannot write auction store")
		return
	}

	_, err = f.file.WriteString(string(bytes))
	if err != nil {
		log.Err(err).Msg("cannot write auction store")
		return
	}

	err = f.file.Sync()
	if err != nil {
		log.Err(err).Msg("cannot sync")
		return
	}
}

func NewFileAuctionRepository() (auction.Repository, error) {
	file, err := os.OpenFile(storePath, os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		return nil, fmt.Errorf("cannnot create scrapper")
	}
	return &FileAuctionRepository{
		file:    file,
		mu:      &sync.Mutex{},
		version: 0,
	}, nil
}
