package commands

import (
	"context"
	"regexp"
)

func ListCommands(router *Router) Command {
	return listCommands{
		router:  router,
		pattern: "^!help",
	}
}

type listCommands struct {
	router  *Router
	pattern string
}

func (c listCommands) Regex() *regexp.Regexp {
	return regexp.MustCompile(c.pattern)
}

func (c listCommands) Execute(ctx context.Context, match []string) (string, error) {
	msg := "Available commands:\n"
	for regex := range c.router.commands {
		msg += "- " + regex.String() + "\n"
	}

	return msg, nil
}
