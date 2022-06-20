// Copyright © 2019 Volodymyr Kalachevskyi <v.kalachevskyi@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cmd is the CLI (command-line interface) for the application.
package cmd

import (
	"fmt"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	h "github.com/Kalachevskyi/mono-chat/app/presetation/telegram"
	"github.com/Kalachevskyi/mono-chat/config"
	"github.com/Kalachevskyi/mono-chat/di"
)

// TimeLocation - application location.
const timeLocation = "Europe/Kiev"

// RootCMD - represents the main command for starting the application.
type RootCMD struct {
	conf config.Config
}

// Init - initializes the CLI application.
func (r *RootCMD) Init() *cli.App {
	cmd := cli.NewApp()
	cmd.Name = "Mono transactions converter"
	cmd.Action = r.serve
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token",
			Usage:       "Telegram token",
			Destination: &r.conf.Token,
			EnvVar:      "TOKEN",
		},
		cli.IntFlag{
			Name:        "offset",
			Usage:       "Telegram offset update",
			Destination: &r.conf.Offset,
			EnvVar:      "OFFSET",
		},
		cli.IntFlag{
			Name:        "timeout",
			Usage:       "Telegram timeout",
			Destination: &r.conf.Timeout,
			EnvVar:      "TIMEOUT",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Telegram debug",
			Destination: &r.conf.Debug,
			EnvVar:      "DEBUG",
		},
		cli.StringFlag{
			Name:        "encoding_log",
			Usage:       `Encoding log, Valid values are "json" and "console"`,
			Destination: &r.conf.EncodingLog,
			Value:       "json",
			EnvVar:      "ENCODING_LOG",
		},
		cli.StringFlag{
			Name:        "redis_url",
			Usage:       `URL to access to redis service"`,
			Destination: &r.conf.RedisURL,
			EnvVar:      "REDIS_URL",
		},
		cli.IntFlag{
			Name:        "http_port",
			Usage:       "Http servier poert",
			Destination: &r.conf.HTTPPort,
			EnvVar:      "HTTP_PORT",
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

	up, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("can't get updates: %v", err.Error())
	}

	log, err := di.Logger(r.conf.Debug, r.conf.EncodingLog)
	if err != nil {
		return err
	}

	rClient, err := di.RedisClient(r.conf.RedisURL)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(timeLocation)
	if err != nil {
		return errors.Wrap(err, "can't set time location")
	}

	toolsWrapper := di.ToolsWrapper{Log: log, RedisClient: rClient, Loc: loc, Bot: bot}
	handlers := map[h.HandlerKey]h.Handler{
		h.FileReportHandler:   di.InjectReport(toolsWrapper),
		h.MappingHandler:      di.InjectMapping(toolsWrapper),
		h.TransactionsHandler: di.InjectTransaction(toolsWrapper),
		h.TokenHandler:        di.InjectToken(toolsWrapper),
		h.ClientInfoHandler:   di.InjectClientInfo(toolsWrapper),
		h.AccountHandler:      di.InjectAccount(toolsWrapper),
		h.ChatUserHandler:     di.InjectUserChat(toolsWrapper),
	}

	fmt.Println("mono_chat_bot is running")

	go h.NewChat(up, handlers, h.NewBotWrapper(bot, log)).Handle()
	httpService := di.InjectHTTPService(toolsWrapper, r.conf.HTTPPort)

	return httpService.Start()
}
