package cmd

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
	"gitlab.com/Kalachevskyi/mono-chat/config"
	"gitlab.com/Kalachevskyi/mono-chat/handlers"
	"gitlab.com/Kalachevskyi/mono-chat/infrastructure"
	"gitlab.com/Kalachevskyi/mono-chat/repository"
	"gitlab.com/Kalachevskyi/mono-chat/usecases"
)

type RootCMD struct {
	conf config.Config
}



func (r *RootCMD) Init() *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "Mono chat converter"
	cmd.Action = r.serve
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token",
			Usage:       "Telegram token",
			Destination: &r.conf.Token,
		},
		cli.IntFlag{
			Name:        "offset",
			Usage:       "Telegram offset update",
			Destination: &r.conf.Offset,
		},
		cli.IntFlag{
			Name:        "timeout",
			Usage:       "Telegram timeout",
			Destination: &r.conf.Timeout,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Telegram debug",
			Destination: &r.conf.Debug,
		},
		cli.StringFlag{
			Name:        "encoding_log",
			Usage:       `Encoding log, Valid values are "json" and "console"`,
			Destination: &r.conf.EncodingLog,
			Value:       "json",
		},
		cli.StringFlag{
			Name:        "mono_api_token",
			Usage:       `Token to access monobank API"`,
			Destination: &r.conf.MonoApiToken,
		},
	}

	return cmd
}






func (r *RootCMD) serve(c *cli.Context) error {
	if err := r.conf.Validate(); err != nil {
		return fmt.Errorf("can't validate config: err=%s", err.Error())
	}

	bot, err := tg.NewBotAPI(r.conf.Token)
	if err != nil {
		return fmt.Errorf("can't initialize Telegram: err=%s", err.Error())
	}

	bot.Debug = r.conf.Debug
	u := tg.NewUpdate(r.conf.Offset)
	u.Timeout = r.conf.Timeout

	zLog, err := infrastructure.GetLogger(r.conf)
	if err != nil {
		return fmt.Errorf("can't load logger: err=%v", err)
	}

	zLog.Infof("Authorized on account %s", bot.Self.UserName)

	dateRegexp := usecases.Date{}
	if err := dateRegexp.Init(); err != nil {
		return fmt.Errorf("can't compaile regexp: err=%s", err.Error())
	}

	chatRepo := repository.NewChat()
	monoRepo := repository.NewMono(r.conf.MonoApiToken, zLog)
	cvsUC := usecases.NewChat(chatRepo, dateRegexp)
	apiUC := usecases.NewApi(monoRepo, dateRegexp)

	chat := handlers.NewChat(bot, u, zLog, cvsUC, apiUC)
	chat.Handle()

	return nil
}
