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

// Package di is an application layer for initializing app components
package di

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/app/infrastructure/mono"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Kalachevskyi/mono-chat/app/infrastructure/redis"
	"github.com/Kalachevskyi/mono-chat/app/infrastructure/telegram"
	h "github.com/Kalachevskyi/mono-chat/app/presetation/api"
	uc "github.com/Kalachevskyi/mono-chat/app/usecases"
	"github.com/Kalachevskyi/mono-chat/config"
)

//Build - initialize app components
func Build(cnf config.Config) (*h.Chat, error) {
	bot, err := tg.NewBotAPI(cnf.Token)
	if err != nil {
		return nil, fmt.Errorf("can't initialize Telegram: err=%s", err.Error())
	}

	bot.Debug = cnf.Debug
	u := tg.NewUpdate(cnf.Offset)
	u.Timeout = cnf.Timeout

	up, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %v", err.Error())
	}

	reportHandler, err := reportHandler(cnf)
	if err != nil {
		return nil, err
	}

	mappingHandler, err := mappingHandler(cnf)
	if err != nil {
		return nil, err
	}

	transactionHandler, err := transactionHandler(cnf)
	if err != nil {
		return nil, err
	}

	tokenHandler, err := tokenHandler(cnf)
	if err != nil {
		return nil, err
	}

	clientInfo, err := clientInfoHandler(cnf)
	if err != nil {
		return nil, err
	}

	accountHandler, err := accountHandler(cnf)
	if err != nil {
		return nil, err
	}

	handlers := map[h.HandlerKey]h.Handler{
		h.FileReportHandler:   reportHandler,
		h.MappingHandler:      mappingHandler,
		h.TransactionsHandler: transactionHandler,
		h.TokenHandler:        tokenHandler,
		h.ClientInfoHandler:   clientInfo,
		h.AccountHandler:      accountHandler,
	}

	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	return h.NewChat(up, handlers, h.NewBotWrapper(tgBot, log)), nil
}

func reportHandler(cnf config.Config) (*h.FileReport, error) {
	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	dateUC, err := uc.NewDate(nil)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	mappingRepo := redis.NewMapping(rClient)
	tgRepo := telegram.NewTelegram(log)
	fileReportUC := uc.NewFileReport(*dateUC, mappingRepo, log, tgRepo)
	fileReportHandler := h.NewFileReport(fileReportUC, h.NewBotWrapper(tgBot, log))
	return fileReportHandler, nil
}

func mappingHandler(cnf config.Config) (*h.Mapping, error) {
	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	mappingRepo := redis.NewMapping(rClient)
	tgRepo := telegram.NewTelegram(log)
	mappingUC := uc.NewMapping(mappingRepo, tgRepo)
	mappingHandler := h.NewMapping(mappingUC, h.NewBotWrapper(tgBot, log))
	return mappingHandler, nil
}

func transactionHandler(cnf config.Config) (*h.Transaction, error) {
	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	dateUC, err := uc.NewDate(nil)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	accountRepo := redis.NewAccount(rClient)
	tokenRepo := redis.NewToken(rClient)
	apiRepo := mono.NewMono(log)
	mappingRepo := redis.NewMapping(rClient)
	// Initialize use cases
	tokenUC := uc.NewToken(tokenRepo)
	accountUC := uc.NewAccount(accountRepo)
	apiUC := uc.NewTransaction(apiRepo, mappingRepo, log, *dateUC)

	transactionHandler := h.NewTransaction(tokenUC, accountUC, apiUC, h.NewBotWrapper(tgBot, log))
	return transactionHandler, nil
}

func tokenHandler(cnf config.Config) (*h.Token, error) {
	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	tokenRepo := redis.NewToken(rClient)
	tokenUC := uc.NewToken(tokenRepo)
	tokenHandler := h.NewToken(tokenUC, h.NewBotWrapper(tgBot, log))
	return tokenHandler, nil
}

func clientInfoHandler(cnf config.Config) (*h.ClientInfo, error) {
	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	tokenRepo := redis.NewToken(rClient)
	clientInfoRepo := mono.NewMono(log)
	tokenUC := uc.NewToken(tokenRepo)
	clientInfoUC := uc.NewClientInfo(clientInfoRepo)
	clientInfo := h.NewClientInfo(tokenUC, clientInfoUC, h.NewBotWrapper(tgBot, log))
	return clientInfo, nil
}

func accountHandler(cnf config.Config) (*h.Account, error) {
	rClient, err := RedisClient(cnf.RedisURL)
	if err != nil {
		return nil, err
	}

	tgBot, err := TGBot(cnf.Token)
	if err != nil {
		return nil, err
	}

	log, err := Logger(cnf.Debug, cnf.EncodingLog)
	if err != nil {
		return nil, err
	}

	accountRepo := redis.NewAccount(rClient)
	accountUC := uc.NewAccount(accountRepo)
	tokenHandler := h.NewAccount(accountUC, h.NewBotWrapper(tgBot, log))
	return tokenHandler, nil
}
