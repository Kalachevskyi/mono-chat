package cmd

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/infrastructure"

	"github.com/Kalachevskyi/mono-chat/config"
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

	chat, err := infrastructure.Build(r.conf)
	if err != nil {
		return err
	}

	fmt.Println("mono_chat_bot is running")

	chat.Handle()

	return nil
}
