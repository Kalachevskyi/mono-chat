package cmd

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/config"
	h "github.com/Kalachevskyi/mono-chat/handlers"
	"github.com/Kalachevskyi/mono-chat/infrastructure"
	"github.com/Kalachevskyi/mono-chat/repository"
	"github.com/Kalachevskyi/mono-chat/usecases"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
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
			Name:        "redis_url",
			Usage:       `URL to access to redis service"`,
			Destination: &r.conf.RedisUrl,
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

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("can't get updates: %v", err.Error())
	}

	redisClient, err := infrastructure.NewRedisClient(r.conf.RedisUrl)
	if err != nil {
		return err
	}

	zLog, err := infrastructure.GetLogger(r.conf)
	if err != nil {
		return fmt.Errorf("can't load logger: err=%v", err)
	}

	dateUC := usecases.Date{}
	if err := dateUC.Init(); err != nil {
		return fmt.Errorf("can't compaile regexp: err=%s", err.Error())
	}

	zLog.Infof("Authorized on account %s", bot.Self.UserName)

	monoRepo := repository.NewMono(zLog)
	tokeRepo := repository.NewToken(redisClient)

	// Initialize usecases
	apiUC := usecases.NewApi(monoRepo, dateUC)
	tokeUC := usecases.NewToken(tokeRepo)

	// Initialize chat handler
	chatBuilder := h.ChatBuilder{
		Updates:  updates,
		ApiUC:    apiUC,
		TokeUC:   tokeUC,
		Handlers: handlers(bot, zLog, dateUC),
	}
	chat := chatBuilder.Build()
	chat.Handle()

	return nil
}
