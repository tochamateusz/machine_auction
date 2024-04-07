package auction

type Repository interface {
	GetAll() []Auction
	Get(id string) Auction
	Save(auction Auction)
}
