package main

import (
	"fmt"
	"os"

	"github.com/zmoog/go-bender/business/bot"
	"github.com/zmoog/go-bender/business/bot/commands"
	"github.com/zmoog/go-bender/foundation/logger"
	"github.com/zmoog/go-bender/foundation/scraper/jsonscraper"
)

var (
	build = "unknown"
	date  = "unknown"
)

func main() {
	// Logging
	log, err := logger.New("bender")
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			log.Errorw("Error syncing logger", "error", err)
		}
	}()

	// Configuration
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		log.Error("DISCORD_TOKEN environment variable not set")
		return
	}

	log.Info("Starting bot")

	// Dependencies
	jsonscraper := jsonscraper.New()

	// Commands
	bender := bot.New(log, token, bot.BuildInfo{
		Version: build,
		Date:    date,
	})
	bender.AddCommand(commands.ListAppleProducts(jsonscraper))

	// Startup
	err = bender.Run()
	if err != nil {
		log.Errorw("Error running bot: %v", err)
	}

	// Shutdown
	log.Info("Bot stopped")
}
