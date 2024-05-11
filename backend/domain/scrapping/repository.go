package scrapping

type Repository interface {
	GetAll() []ScrappedAuctions
	Get(id string) []ScrappedAuctions
	Save(auction ScrappedAuctions)
}
