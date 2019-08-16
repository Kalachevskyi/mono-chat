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

// Package cmd is the CLI (command-line interface) for the application
package cmd

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/config"
	"github.com/Kalachevskyi/mono-chat/infrastructure"
	"github.com/urfave/cli"
)

// RootCMD - represents the main command for starting the application
type RootCMD struct {
	conf config.Config
}

// Init - initializes the CLI application
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
			Destination: &r.conf.RedisURL,
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
