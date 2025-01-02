package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/zmoog/go-bender/foundation/scraper/jsonscraper"
)

func ListAppleProducts(scraper jsonscraper.Scraper) Command {
	return listAppleProducts{
		pattern: "^!apple (ca|cn|fr|it|uk|us) (accessories|airpods|appletvs|homepods|iphones|ipads|macs)(\\s+.*)?",
		scraper: scraper,
	}
}

type listAppleProducts struct {
	pattern string
	scraper jsonscraper.Scraper
}

func (c listAppleProducts) Regex() *regexp.Regexp {
	return regexp.MustCompile(c.pattern)
}

func (c listAppleProducts) Execute(ctx context.Context, match []string) (string, error) {
	var products []product
	var filteredProducts []product

	url := fmt.Sprintf("https://raw.githubusercontent.com/zmoog/refurbished-history/main/stores/%s/%s.json", match[1], match[2])
	err := c.scraper.Scrape(url, &products)
	if err != nil {
		return "", err
	}

	for _, p := range products {
		if len(match) == 4 && !strings.Contains(p.Name, match[3]) {
			continue
		}
		filteredProducts = append(filteredProducts, p)
	}

	if len(filteredProducts) == 0 {
		return fmt.Sprintf("No %s in stock in %s", match[2], match[1]), nil
	}

	msg := "Products in stock:\n"
	for _, p := range filteredProducts {
		msg += fmt.Sprintf("- %s %v\n", p.Name, p.Price)
	}

	return msg, err
}

type product struct {
	Name  string  `json:"name"`
	URL   string  `json:"url"`
	Price float64 `json:"price"`
}
