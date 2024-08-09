package scrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/goutil/dump"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

const domain = "https://www.ab-auction.com/pl"
const myObservation = "https://www.ab-auction.com/pl/account/myobservations"

const scrappingResultDir = "./scrapping-result/"

type DescriptionFounded struct {
	Id          string
	Description []string
}

type StartingPriceFound struct {
	Id            string
	StartingPrice string
}

type Scrapper struct {
	client     *http.Client
	file       *os.File
	repository auction.Repository

	done               chan struct{ Id string }
	descriptionFounded chan DescriptionFounded
	startingPriceFound chan StartingPriceFound

	auctionsMap map[string]*Fullfill

	mu *sync.Mutex
}

func NewScrapper(repository auction.Repository) (*Scrapper, error) {

	indexHtml := scrappingResultDir + "index.html"
	fileInfo, err := os.Stat(indexHtml)

	log.Info().Msgf("Checking if file exist:")
	var file *os.File
	if errors.Is(err, os.ErrNotExist) {
		log.Info().Msgf("Creating")
		file, err = os.Create(indexHtml)
		if err != nil {
			log.Err(err).Msgf("bad request")
			return nil, errors.Join(err, fmt.Errorf("cannot create scrapper"))
		}
		fileInfo, err = file.Stat()
		if err != nil {

			log.Err(err).Msgf("bad request")
			return nil, errors.Join(err, fmt.Errorf("cannot create scrapper"))
		}
	} else {
		file, err = os.OpenFile(indexHtml, os.O_RDWR, 0777)
		if err != nil {
			log.Err(err).Msgf("bad request")
			return nil, errors.Join(err, fmt.Errorf("cannot create scrapper"))
		}
	}
	log.Info().Msgf("Founded: %+v\n", fileInfo.Name())

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}
	scrapper := &Scrapper{
		client,
		file,
		repository,
		make(chan struct{ Id string }),
		make(chan DescriptionFounded),
		make(chan StartingPriceFound),
		nil,
		&sync.Mutex{},
	}

	go scrapper.listen(context.Background())
	return scrapper, nil
}

func (s *Scrapper) Scrap(ctx context.Context, auctions []auction.Auction) map[string]*Fullfill {

	mapAuctions := make(map[string]*Fullfill)

	mapAuctions = lo.Reduce(auctions,
		func(agg map[string]*Fullfill, a auction.Auction, _ int) map[string]*Fullfill {
			agg[a.Id()] = NewFullfillPorcess(a)
			return agg
		}, mapAuctions)

	s.auctionsMap = mapAuctions

	return mapAuctions
}

func (s *Scrapper) listen(_ context.Context) {
listenLoop:
	for {
		select {
		case done := <-s.done:
			{
				dump.P(s.auctionsMap)
				process, ok := s.auctionsMap[done.Id]
				if !ok {
					log.Logger.Debug().Caller().Err(errors.New("auction id: [" + done.Id + "] not exist")).Msg("")
					continue
				}
				process.Done()
				delete(s.auctionsMap, done.Id)
				s.repository.Save(process.rawAuction)
				log.Info().Msgf("Auction Id:%s is done. %d left", done.Id, len(s.auctionsMap))
				if len(s.auctionsMap) <= 0 {
					break listenLoop
				}
			}

		case description := <-s.descriptionFounded:
			{
				process, ok := s.auctionsMap[description.Id]
				if !ok {
					log.Logger.Debug().Caller().Err(errors.New("auction id: [" + description.Id + "] not exist")).Msg("")
					continue
				}
				process.Description(description.Description)
			}

		case startingPrice := <-s.startingPriceFound:
			{
				process, ok := s.auctionsMap[startingPrice.Id]
				if !ok {
					log.Logger.Debug().Caller().Err(errors.New("auction id: [" + startingPrice.Id + "] not exist")).Msg("")
					continue
				}
				process.StartingPrice(startingPrice.StartingPrice)
			}
		}
	}
}

func (s *Scrapper) OnAuctionFound(ctx context.Context, message interface{}) {
	auctionFounded, ok := message.(*Fullfill)
	if ok == false {
		log.Err(fmt.Errorf("can't parse auction found message")).Msgf("")
	}

	log.Info().
		Str("AuctionId", auctionFounded.rawAuction.Id()).
		Str("AuctionName", auctionFounded.rawAuction.Name()).
		Msg("Auction requesting...")

	req, err := http.NewRequest("GET", domain+"/auction/"+auctionFounded.rawAuction.Id(), nil)
	if err != nil {
		log.Err(err).Msgf("can't get auction id: %s", auctionFounded.rawAuction.Id())
	}

	res, err := s.client.Do(req)
	if err != nil {
		log.Err(err).Msgf("bad request")
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Debug().Caller().Err(err).Msgf("Incorrect status code: %+v\n", res.Status)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Msgf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err).Msgf("cant read html")
	}

	err = os.MkdirAll("./scrapping-result/"+auctionFounded.rawAuction.Id()+"/", 0777)
	if err != nil {
		panic(err)
	}

	f, e := os.Create("./scrapping-result/" + auctionFounded.rawAuction.Id() + "/index.html") // "m1UIjW1.jpg"
	if e != nil {
		panic(e)
	}
	defer f.Close()
	html, _ := doc.Html()
	f.WriteString(html)

	selection := doc.Find(".swiper-wrapper")
	selection.Find(".swiper-slide").Each(func(i int, sel *goquery.Selection) {
		imageSrc, exist := sel.Find(".img-fluid").Attr("src")
		if exist == true {
			log.Info().Msgf("Image source: %s", imageSrc)
			s.SaveImage(imageSrc, "./scrapping-result/"+auctionFounded.rawAuction.Id()+"/"+fmt.Sprintf("%d", i)+".jpg")
		}
	})

	detailSelection := doc.Find(".details")
	description := []string{}
	detailSelection.Each(func(i int, s *goquery.Selection) {
		_ = s.Find(".row").Each(func(i int, divSel *goquery.Selection) {
			description = append(description, divSel.Text())
		})
	})

	detailFile, e := os.Create("./scrapping-result/" + auctionFounded.rawAuction.Id() + "/detail.html") // "m1UIjW1.jpg"
	htmlDetailSection, _ := detailSelection.Html()

	detailFile.WriteString(htmlDetailSection)

	startingPrice := strings.TrimSpace(doc.Find("div.mt-n2:nth-child(1) > span:nth-child(1)").Text())

	s.descriptionFounded <- DescriptionFounded{
		Id:          auctionFounded.rawAuction.Id(),
		Description: description,
	}

	s.startingPriceFound <- StartingPriceFound{
		Id:            auctionFounded.rawAuction.Id(),
		StartingPrice: startingPrice,
	}

	s.done <- struct{ Id string }{
		Id: auctionFounded.rawAuction.Id(),
	}

}

func (s *Scrapper) SaveAuctionAssets(a auction.Auction) {

}

func (s *Scrapper) SaveImage(url string, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := os.Stat(path)
	if err == nil && !errors.Is(err, os.ErrNotExist) {
		log.Info().Msgf("File exist: %+v\n", path)
		return
	}
	f, err := os.Create(path)
	if err != nil {
		log.Err(err).Msgf("Cannot create  file: %+v\n", path)
		return
	}
	defer f.Close()

	client := &http.Client{
		Timeout: time.Duration(time.Second * 50),
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msgf("new request failed  %+v\n", url)
		return
	}
	r, err := client.Do(req)
	if err != nil {
		log.Err(err).Msgf("cannot get image  %+v\n", url)
		return
	}
	defer r.Body.Close()

	n, err := f.ReadFrom(r.Body)
	if err != nil {
		log.Err(err).Msgf("cannot read body  %+v\n", n)
		return
	}
	fmt.Println("File size: ", n)
}

func (s *Scrapper) GetAuctions() ([]auction.Auction, error) {
	s.Login()
	auctions, err := s.getMyObservation()
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (s *Scrapper) getCookie() ([]*http.Cookie, error) {
	domainUrl, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	cookies := s.client.Jar.Cookies(domainUrl)
	return cookies, nil
}

func (s *Scrapper) getMyObservation() ([]auction.Auction, error) {
	req, err := http.NewRequest("GET", myObservation, nil)
	if err != nil {
		log.Debug().Caller().Err(err).Msgf("Bade new request: %+v\n", req.URL)
		return nil, err
	}

	res, err := s.client.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Err(err).Msgf("Requests to: %s Incorrect status code: %+v\n", req.URL, res.Status)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error().Msgf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err).Msgf("cant read html")
	}

	selection := doc.Find("div.container:nth-child(6) > table:nth-child(2) > tbody:nth-child(2)")
	var auctions []auction.Auction
	selection.Find("tr").Each(func(i int, s *goquery.Selection) {
		id := s.Find("td:nth-child(1)").Text()
		image := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
		name := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		year := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		price := strings.TrimSpace(s.Find("td:nth-child(5)").Text())
		endDate := strings.TrimSpace(s.Find("td:nth-child(6)").Text())

		auction := auction.NewAuction(id, image, name, year, price, endDate)
		auctions = append(auctions, *auction)
	})

	html, err := selection.Html()
	if err != nil {
		return nil, err
	}

	s.file.WriteString("<htmL><body><table>" + html + "</table></body></html>")
	return auctions, nil
}

func (s *Scrapper) PrintCookie() {
	domainUrl, err := url.Parse(domain)
	if err != nil {
		log.Err(err).Msgf("cannot get url of domain: %s", domain)
		return
	}
	for _, v := range s.client.Jar.Cookies(domainUrl) {
		log.Info().Msgf("Cookie: %+v\n", v)
	}
}

func (s *Scrapper) GetAllImages(id string) ([]string, error) {
	files, err := os.ReadDir(scrappingResultDir + "/" + id)
	if err != nil {
		return nil, err
	}

	images := []string{}
	for _, f := range files {
		var validImage = regexp.MustCompile(`^(.*).jpg$`)
		if validImage.Match([]byte(f.Name())) {
			images = append(images, f.Name())
		}
	}
	return images, nil
}

func (s *Scrapper) Login() error {

	params := url.Values{}

	LOGIN := os.Getenv("LOGIN")
	params.Add("login_user", LOGIN)

	PASSWORD := os.Getenv("PASSWORD")
	params.Add("login_passwort", PASSWORD)
	params.Add("login", "login")

	postData := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", domain+"/login", postData)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.Do(req)
	if err != nil {
		log.Err(err).Msgf("%v can't login", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Err(errors.New("can't login")).Msgf("")
		return err
	}

	s.PrintCookie()

	log.Info().Msgf("Succesfull loged in")

	return nil
}
