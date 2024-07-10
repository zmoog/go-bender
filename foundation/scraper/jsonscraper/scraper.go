package jsonscraper

import (
	"encoding/json"
	"net/http"
)

type Scraper interface {
	Scrape(url string, v interface{}) error
}

func New() Scraper {
	return scraper{}
}

type scraper struct{}

func (s scraper) Scrape(url string, v interface{}) error {
	// HTTP get the JSON document in the url and unmarshall it into v

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
