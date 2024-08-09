package scrapper

import "github.com/tochamateusz/machine_auction/domain/auction"

type fullfillState interface {
	Transit(event any) fullfillState
}

type Fullfill struct {
	state      string
	rawAuction auction.Auction
	Name       string
}

func NewFullfillPorcess(a auction.Auction) *Fullfill {
	return &Fullfill{
		state:      "init",
		rawAuction: a,
		Name:       a.Name(),
	}
}

func (f *Fullfill) StartingPrice(s string) {
	f.rawAuction.DefineStartingPrice(s)
}

func (f *Fullfill) Description(description []string) {
	f.rawAuction.Describe(description)
}

func (f *Fullfill) Done() auction.Auction {
	f.state = "done"
	return f.rawAuction
}
