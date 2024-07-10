package main

import (
	"fmt"
	"os"

	"github.com/zmoog/go-bender/business/bot"
	"github.com/zmoog/go-bender/business/bot/commands"
	"github.com/zmoog/go-bender/foundation/logger"
	"github.com/zmoog/go-bender/foundation/scraper/jsonscraper"
)

func main() {
	// Logging
	log, err := logger.New("bender")
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

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
	bender := bot.New(log, token)
	bender.AddCommand(commands.ListAppleProducts(jsonscraper))

	// Startup
	err = bender.Run()
	if err != nil {
		log.Errorw("Error running bot: %v", err)
	}

	// Shutdown
	log.Info("Bot stopped")
}
