package cmd

import (
	h "github.com/Kalachevskyi/mono-chat/handlers"
	"github.com/Kalachevskyi/mono-chat/repository"
	uc "github.com/Kalachevskyi/mono-chat/usecases"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handlers(bot *tg.BotAPI, log h.Logger, dateUC uc.Date) map[h.HandlerKey]h.Handler {
	return map[h.HandlerKey]h.Handler{
		h.FileReportHandler: reportHandler(bot, log, dateUC),
		h.MappingHandler:    mappingHandler(bot, log),
	}
}

func reportHandler(bot *tg.BotAPI, log h.Logger, dateUC uc.Date) *h.FileReport {
	repo := repository.NewTelegram()
	fileReportUC := uc.NewFileReport(repo, dateUC)
	botWrapper := h.NewBotWrapper(bot, log)
	fileReportHandler := h.NewFileReport(fileReportUC, botWrapper)
	return fileReportHandler
}

func mappingHandler(bot *tg.BotAPI, log h.Logger) *h.Mapping {
	repo := repository.NewTelegram()
	mappingUC := uc.NewMapping(repo)
	botWrapper := h.NewBotWrapper(bot, log)
	mappingHandler := h.NewMapping(mappingUC, botWrapper)
	return mappingHandler
}

func transactionHandler(bot *tg.BotAPI, log h.Logger) *h.Transaction {
	tokenUC := uc.NewToken()
	return nil
}
