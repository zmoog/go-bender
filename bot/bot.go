package bot

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	log   *zap.SugaredLogger
	token string
}

func New(log *zap.SugaredLogger, token string) *Bot {
	return &Bot{
		log:   log,
		token: token,
	}
}

func (b *Bot) Run() error {
	// Override the default logging function to use our logger.
	discordgo.Logger = b.logger

	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	session.AddHandler(b.newMessage)

	err = session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection to Discord: %w", err)
	}
	defer session.Close()

	b.log.Info("Bot is now running. Press CTRL-C to exit.")

	// Wait here until CTRL-C or other term signal is received.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	b.log.Info("Shutting down bot.")

	return nil
}

func (b *Bot) newMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	b.log.Infow(
		"Received message",
		"message", m.Content,
		"author", m.Author.Username,
	)

	switch {
	case m.Content == "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			b.log.Error("Error sending message: ", err)
		}
	}
}

func (b *Bot) logger(msgL, caller int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	// Get the caller's info (file, line, function)
	originalCaller := "unknown:0"
	_, file, line, ok := runtime.Caller(caller + 1) // Skip 1 caller to get the original caller
	if ok {
		originalCaller = fmt.Sprintf("%s:%d", file, line)
	}

	switch msgL {
	case discordgo.LogError:
		b.log.Errorw(msg, "original_caller", originalCaller)
	case discordgo.LogWarning:
		b.log.Warnw(msg, "original_caller", originalCaller)
	case discordgo.LogInformational:
		b.log.Infow(msg, "original_caller", originalCaller)
	case discordgo.LogDebug:
		b.log.Debugw(msg, "original_caller", originalCaller)
	}
}
