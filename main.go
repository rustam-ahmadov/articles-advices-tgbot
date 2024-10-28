package main

import (
	"flag"
	"log"
	"m/client/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

}

func mustToken() string {
	token := flag.String("bot-token", "", "tg-bot access token")

	flag.Parse()

	if *token == "" {
		log.Fatalf("token has not been provided")
	}

	return *token
}
