package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/zmoog/go-bender/scraper/jsonscraper"
)

func ListAppleProducts(scraper jsonscraper.Scraper) Command {
	return listAppleProducts{
		pattern: "!apple (ca|cn|fr|it|uk|us) (iphones|ipads|macs)(\\s+.*)?",
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
	var products []struct {
		Name  string  `json:"name"`
		URL   string  `json:"url"`
		Price float64 `json:"price"`
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/zmoog/refurbished-history/main/stores/%s/%s.json", match[1], match[2])
	err := c.scraper.Scrape(url, &products)
	if err != nil {
		return "", err
	}

	// TODO: use Go templates to format the message instead of
	// concatenating strings like an animal.
	msg := "Products in stock:\n"
	for _, p := range products {
		if len(match) == 4 && !strings.Contains(p.Name, match[3]) {
			continue
		}
		msg += fmt.Sprintf("- %s %v\n", p.Name, p.Price)
	}

	return msg, err
}
