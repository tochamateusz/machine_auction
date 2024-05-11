package auction

type AuctionFounded struct{ Auction Auction }
type AuctionsFounded struct {
	Id      string
	Auction []Auction
}
type AuctionScrapped struct{}
type AuctionSaved struct{}
