package scrapper

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

const domain = "https://www.ab-auction.com"
const myObservation = "https://www.ab-auction.com/pl/account/myobservations"

type Scrapper struct {
	client *http.Client
	file   *os.File
}

func NewScrapper() (*Scrapper, error) {
	file, err := os.Create("./web/build/index.html")
	if err != nil {
		return nil, fmt.Errorf("cannnot create scrapper")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	return &Scrapper{
		client,
		file,
	}, nil
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
		return nil, err
	}

	res, err := s.client.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Err(err).Msgf("Incorrect status code: %+v\n", res.Status)
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

func (s *Scrapper) Login() error {

	params := url.Values{}

	LOGIN := os.Getenv("LOGIN")
	params.Add("login_user", LOGIN)

	PASSWORD := os.Getenv("PASSWORD")
	params.Add("login_passwort", PASSWORD)
	params.Add("login", "login")

	postData := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", domain+"/pl/login", postData)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	_, err = s.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
