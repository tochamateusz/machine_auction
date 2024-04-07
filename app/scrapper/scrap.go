package scrapper

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const domain = "https://www.ab-auction.com"

type Scrapper struct {
	client *http.Client
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
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	s.client = &http.Client{
		Jar: jar,
	}

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

	_, err = s.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
