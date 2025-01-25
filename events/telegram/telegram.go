package telegram

import "main/clients/telegram"

type EventProcessor struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(client *telegram.Client) {

}
