package domain

import "github.com/tochamateusz/machine_auction/domain/auction"

type fullfillState interface {
	Transit(event any) fullfillState
}

type Fullfill struct {
	fullfillState
	rawAuction auction.Auction
	Name       string
}

func NewFullfillPorcess(a auction.Auction) *Fullfill {
	return &Fullfill{
		rawAuction: a,
		Name:       a.Name(),
	}
}

func (f *Fullfill) IsDone() bool {
	return false
}

func (f *Fullfill) Done() {
	f.fullfillState.Transit(struct {
		name string
	}{
		name: "done",
	})
}
