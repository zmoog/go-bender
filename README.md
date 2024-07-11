# go-bender

## how to run 

### requirements

- set the `DISCORD_TOKEN` environment variable

```shell
export DISCORD_TOKEN="<redacted>" 
```

### locally

```shell
$ make run-local
go run app/services/bender-bot/main.go | go run app/tooling/logfmt/main.go
bender: 2024-06-01T07:27:44.733+0200: bender-bot/main.go:28: info: 00000000-0000-0000-0000-000000000000: Starting bot
bender: 2024-06-01T07:27:45.791+0200: bot/bot.go:49: info: 00000000-0000-0000-0000-000000000000: Bot is now running. Press CTRL-C to exit.
```

### on k8s

1. start a kind cluster