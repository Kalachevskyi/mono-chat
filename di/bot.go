package di

import (
	"fmt"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	tgBot  *tg.BotAPI //nolint:gochecknoglobals
	tgOnce sync.Once  //nolint:gochecknoglobals
)

// TGBot - initialize the instance of Telegram client
func TGBot(token string) (*tg.BotAPI, error) {
	var err error

	tgOnce.Do(func() {
		tgBot, err = tg.NewBotAPI(token)
		if err != nil {
			err = fmt.Errorf("can't initialize Telegram client: err=%s", err.Error())
			return
		}
	})

	return tgBot, err
}
