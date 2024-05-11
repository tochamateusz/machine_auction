package scrapping

import "github.com/tochamateusz/machine_auction/domain/auction"

type ScrappedAuctions struct {
	Id        string
	Triggered string
	Auction   auction.Auction
}
