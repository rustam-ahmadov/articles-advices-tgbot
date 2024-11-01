package telegram

import "m/client/telegram"

type TgProcessor struct {
	tg     *telegram.Client
	offset int
	//storage
}

func New(client *telegram.Client, storage)
