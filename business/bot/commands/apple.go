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
	var allProducts []product
	var matchingProducts []product

	url := fmt.Sprintf("https://raw.githubusercontent.com/zmoog/refurbished-history/main/stores/%s/%s.json", match[1], match[2])
	err := c.scraper.Scrape(url, &allProducts)
	if err != nil {
		return "", err
	}

	if len(allProducts) == 0 {
		return fmt.Sprintf("No %s in stock in %s", match[2], match[1]), nil
	}

	// We have at least one product in stock. Not we need to filter
	// the products based on the third argument if it exists.
	switch len(match) {
	case 3: // No filter
		matchingProducts = allProducts
	case 4: // Filter by third argument
		for _, p := range allProducts {
			if !strings.Contains(p.Name, match[3]) {
				continue
			}
			matchingProducts = append(matchingProducts, p)
		}

		// Check if we have any matching products.
		if len(matchingProducts) == 0 {
			return fmt.Sprintf(
				"No %s with %s in stock in %s",
				match[2],
				strings.TrimSpace(match[3]),
				match[1],
			), nil
		}
	}

	// We have at least one matching product.
	// Render the message with the matching
	// products.
	msg := "Products in stock:\n"
	for _, p := range matchingProducts {
		msg += fmt.Sprintf("- %s %v\n", p.Name, p.Price)
	}

	return msg, err
}

type product struct {
	Name  string  `json:"name"`
	URL   string  `json:"url"`
	Price float64 `json:"price"`
}
