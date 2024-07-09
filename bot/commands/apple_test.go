package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAppleProducts(t *testing.T) {
	scraper := MockScraper{}
	scraper.On("https://raw.githubusercontent.com/zmoog/refurbished-history/main/stores/it/iphones.json", `[
  {
	"name": "iPhone 12 Pro 256GB",
	"url": "https://apple.com/iphone-12-pro-max-256gb",
	"price": 999.99
  },
  {
	"name": "iPhone 12 Pro Max 256GB",
	"url": "https://apple.com/iphone-12-pro-max-256gb",
	"price": 1099.99
  }
]
`)

	router := NewRouter()
	router.Register(ListAppleProducts(&scraper))

	t.Run("SimpleCommand", func(t *testing.T) {
		found, msg, err := router.FindAndExecute("!apple it iphones")
		require.NoError(t, err)

		assert.True(t, found)
		assert.Equal(t, `Products in stock:
- iPhone 12 Pro 256GB 999.99
- iPhone 12 Pro Max 256GB 1099.99
`, msg)
	})

	t.Run("CommandWithNameFilter", func(t *testing.T) {
		found, msg, err := router.FindAndExecute("!apple it iphones Pro Max")
		require.NoError(t, err)

		assert.True(t, found)
		assert.Equal(t, `Products in stock:
- iPhone 12 Pro Max 256GB 1099.99
`, msg)
	})

	t.Run("NonExistingProduct", func(t *testing.T) {
		found, _, err := router.FindAndExecute("!apple it screwdriwer")
		assert.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("NonExistingCountry", func(t *testing.T) {
		found, _, err := router.FindAndExecute("!apple ie iphones")
		assert.NoError(t, err)
		assert.False(t, found)
	})
}
