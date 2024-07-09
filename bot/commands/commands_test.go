package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type MockScraper struct {
	responses map[string]string
}

func (m *MockScraper) On(url, response string) {
	if m.responses == nil {
		m.responses = make(map[string]string)
	}
	m.responses[url] = response
}

func (m *MockScraper) Scrape(url string, v interface{}) error {
	response, exists := m.responses[url]
	if !exists {
		return fmt.Errorf("no response for url %s", url)
	}

	err := json.NewDecoder(bytes.NewBufferString(response)).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
