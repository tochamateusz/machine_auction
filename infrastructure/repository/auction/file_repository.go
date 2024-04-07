package auction_file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

const storePath = "./db/auctions.json"

type FileAuctionRepository struct {
	file *os.File
}

// Get implements AuctionRepository.
func (*FileAuctionRepository) Get(id string) auction.Auction {
	raw, _ := os.ReadFile(storePath)
	auctionDataModel := make(map[string]AuctionDataModel)

	json.Unmarshal(raw, &auctionDataModel)
	v, ok := auctionDataModel[id]
	if ok == false {
		return auction.Auction{}
	}

	return *auction.NewAuction(v.Id, v.Image, v.Name, v.Year, v.Price, v.EndDate)
}

// GetAll implements AuctionRepository.
func (*FileAuctionRepository) GetAll() []auction.Auction {
	raw, _ := os.ReadFile(storePath)
	auctionDataModel := make(map[string]AuctionDataModel)

	json.Unmarshal(raw, &auctionDataModel)

	var array []auction.Auction
	for _, v := range auctionDataModel {
		array = append(array, *auction.NewAuction(v.Id, v.Image, v.Name, v.Year, v.Price, v.EndDate))
	}
	return array
}

type AuctionDataModel struct {
	Id      string `json:"id"`
	Image   string `json:"image"`
	Name    string `json:"name"`
	Year    string `json:"year"`
	Price   string `json:"price"`
	EndDate string `json:"end_date"`
}

// Save implements AuctionRepository.
func (f *FileAuctionRepository) Save(a auction.Auction) {
	raw, _ := os.ReadFile(storePath)
	auctionDataModel := make(map[string]AuctionDataModel)

	json.Unmarshal(raw, &auctionDataModel)

	auctionDataModel[a.Id()] = AuctionDataModel{
		Id:      a.Id(),
		Image:   a.Image(),
		Name:    a.Name(),
		Year:    a.Year(),
		Price:   a.Price(),
		EndDate: a.EndDate(),
	}

	bytes, err := json.Marshal(auctionDataModel)
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

func NewFileAuctionRepository() (auction.Repository, error) {
	file, err := os.OpenFile(storePath, os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		return nil, fmt.Errorf("cannnot create scrapper")
	}
	return &FileAuctionRepository{
		file,
	}, nil
}
