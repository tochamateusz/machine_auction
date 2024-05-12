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

	// "time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

const domain = "https://www.ab-auction.com"
const myObservation = "https://www.ab-auction.com/pl/account/myobservations"

const scrappingResultDir = "./scrapping-result/"

type DescriptionFounded struct {
	Id          string
	Description []string
}

type Scrapper struct {
	client             *http.Client
	file               *os.File
	done               chan struct{ Id string }
	descriptionFounded chan DescriptionFounded

	onDone             func(id string)
	onDescriptionFound func(DescriptionFounded)

	mu *sync.Mutex
}

func NewScrapper() (*Scrapper, error) {

	indexHtml := scrappingResultDir + "index.html"
	fileInfo, err := os.Stat(indexHtml)

	log.Info().Msgf("Checking if file exist:")
	var file *os.File
	if errors.Is(err, os.ErrNotExist) {
		log.Info().Msgf("Creating")
		file, err = os.Create(indexHtml)
		if err != nil {

			log.Err(err).Msgf("bad request")
			return nil, fmt.Errorf("cannnot create scrapper")
		}
		fileInfo, err = file.Stat()
		if err != nil {

			log.Err(err).Msgf("bad request")
			return nil, fmt.Errorf("cannnot create scrapper")
		}
	} else {
		file, err = os.OpenFile(indexHtml, os.O_RDWR, 0777)
		if err != nil {
			log.Err(err).Msgf("bad request")
			return nil, fmt.Errorf("cannnot create scrapper")
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
		make(chan struct{ Id string }),
		make(chan DescriptionFounded),
		func(id string) {},
		func(DescriptionFounded) {},
		&sync.Mutex{},
	}

	go scrapper.listen(context.Background())

	return scrapper, nil
}

func (s *Scrapper) listen(ctx context.Context) {
	for {
		select {
		case done := <-s.done:
			{
				s.onDone(done.Id)
			}

		case description := <-s.descriptionFounded:
			{
				s.onDescriptionFound(description)
			}
		}
	}
}

func (s *Scrapper) RegisterDone(onDone func(done string)) {
	s.onDone = onDone
}

func (s *Scrapper) RegisterOnDescriptionFound(onDescriptionFound func(DescriptionFounded)) {
	s.onDescriptionFound = onDescriptionFound
}

func (s *Scrapper) OnAuctionFound(ctx context.Context, message interface{}) {
	auctionFounded, ok := message.(auction.AuctionFounded)
	if ok == false {
		log.Err(fmt.Errorf("can't parse auction found message")).Msgf("")
	}

	log.Info().
		Str("AuctionId", auctionFounded.Auction.Id()).
		Str("AuctionName", auctionFounded.Auction.Name()).
		Msg("Auction requesting...")

	req, err := http.NewRequest("GET", domain+"/auction/"+auctionFounded.Auction.Id(), nil)
	if err != nil {
		log.Err(err).Msgf("can't get auction id: %s", auctionFounded.Auction.Id())
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

	err = os.MkdirAll("./scrapping-result/"+auctionFounded.Auction.Id()+"/", 0777)
	if err != nil {
		panic(err)
	}

	f, e := os.Create("./scrapping-result/" + auctionFounded.Auction.Id() + "/index.html") // "m1UIjW1.jpg"
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
			s.SaveImage(imageSrc, "./scrapping-result/"+auctionFounded.Auction.Id()+"/"+fmt.Sprintf("%d", i)+".jpg")
		}
	})

	detailSelection := doc.Find(".details")
	description := []string{}
	detailSelection.Each(func(i int, s *goquery.Selection) {
		_ = s.Find(".row").Each(func(i int, divSel *goquery.Selection) {
			description = append(description, divSel.Text())
		})
	})

	s.descriptionFounded <- DescriptionFounded{
		Id:          auctionFounded.Auction.Id(),
		Description: description,
	}

	detailFile, e := os.Create("./scrapping-result/" + auctionFounded.Auction.Id() + "/detail.html") // "m1UIjW1.jpg"
	htmlDetailSection, _ := detailSelection.Html()

	detailFile.WriteString(htmlDetailSection)

	s.done <- struct{ Id string }{
		Id: auctionFounded.Auction.Id(),
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
		log.Err(err).Msgf("cannot get image  %+v\n", r.Request.URL)
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
	req, err := http.NewRequest("POST", domain+"/pl/login", postData)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = s.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
