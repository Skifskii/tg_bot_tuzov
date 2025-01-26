package telegram

import (
	"errors"
	"log"
	"main/lib/e"
	"main/storage"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (ep *EventProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return ep.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return ep.sendRandom(chatID, username)
	case HelpCmd:
		return ep.sendHelp(chatID)
	case StartCmd:
		return ep.sendHello(chatID)
	default:
		return ep.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (ep *EventProcessor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := ep.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return ep.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := ep.storage.Save(page); err != nil {
		return err
	}

	if err := ep.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (ep *EventProcessor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

	page, err := ep.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return ep.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := ep.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return ep.storage.Remove(page)
}

func (ep *EventProcessor) sendHelp(chatID int) error {
	return ep.tg.SendMessage(chatID, msgHelp)
}

func (ep *EventProcessor) sendHello(chatID int) error {
	return ep.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
