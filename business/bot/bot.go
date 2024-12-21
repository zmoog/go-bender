package bot

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
	"github.com/zmoog/go-bender/business/bot/commands"
)

type BuildInfo struct {
	Version string
	Date    string
}

type Bot struct {
	log       *zap.SugaredLogger
	token     string
	router    commands.Router
	buildInfo BuildInfo
}

func New(log *zap.SugaredLogger, token string, buildInfo BuildInfo) *Bot {
	router := commands.NewRouter()

	return &Bot{
		log:       log,
		token:     token,
		router:    router,
		buildInfo: buildInfo,
	}
}

func (b *Bot) AddCommand(c commands.Command) {
	b.router.Register(c)
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

	b.log.Info("Bot is now running. Press CTRL-C to exit.", "version", b.buildInfo.Version, "date", b.buildInfo.Date)

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
		"channel_id", m.ChannelID,
		"type", m.Type,
	)

	found, msg, err := b.router.FindAndExecute(m.Content)
	if err != nil {
		b.log.Error("Error executing command: ", err)
		return
	}

	if !found {
		b.log.Infof("no command for: %s", m.Content)
		return
	}

	msg = msg[:min(2000, len(msg))]

	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		b.log.Error("Error sending message: ", err)
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
