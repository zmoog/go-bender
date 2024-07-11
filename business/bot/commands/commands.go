package commands

import (
	"context"
	"regexp"
)

type Command interface {
	Regex() *regexp.Regexp
	Execute(ctx context.Context, match []string) (string, error)
}
