package auction

type Auction struct {
	id      string
	image   string
	name    string
	year    string
	price   string
	endDate string
}

func NewAuction(id, image, name, year, price, endDate string) *Auction {
	return &Auction{
		id, image, name, year, price, endDate,
	}
}

func (a *Auction) Id() string {
	return a.id
}

func (a *Auction) Image() string {
	return a.image
}

func (a *Auction) Name() string {
	return a.name
}

func (a *Auction) Year() string {
	return a.year
}

func (a *Auction) Price() string {
	return a.price
}

func (a *Auction) EndDate() string {
	return a.endDate
}
