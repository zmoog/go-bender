package commands

import (
	"context"
	"regexp"
)

func NewRouter() Router {
	router := Router{
		commands: make(map[*regexp.Regexp]Command),
	}

	router.Register(ListCommands(&router))

	return router
}

type Router struct {
	commands map[*regexp.Regexp]Command
}

func (r Router) Register(cmd Command) {
	r.commands[cmd.Regex()] = cmd
}

func (r Router) FindAndExecute(msg string) (bool, string, error) {
	for regex, cmd := range r.commands {
		match := regex.FindStringSubmatch(msg)
		if match != nil {
			msg, err := cmd.Execute(context.Background(), match)
			return true, msg, err
		}
	}

	return false, "", nil
}
