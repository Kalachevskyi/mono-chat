package infrastructure

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/config"
	h "github.com/Kalachevskyi/mono-chat/handlers"
	"github.com/Kalachevskyi/mono-chat/repository"
	uc "github.com/Kalachevskyi/mono-chat/usecases"
	"github.com/go-redis/redis"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

var (
	log         *zap.SugaredLogger
	redisClient *redis.Client
	botWrapper  *h.BotWrapper
	dateUC      uc.Date
)

func Build(conf config.Config) (*h.Chat, error) {
	bot, err := tg.NewBotAPI(conf.Token)
	if err != nil {
		return nil, fmt.Errorf("can't initialize Telegram: err=%s", err.Error())
	}

	bot.Debug = conf.Debug
	u := tg.NewUpdate(conf.Offset)
	u.Timeout = conf.Timeout

	up, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %v", err.Error())
	}

	rClient, err := NewRedisClient(conf.RedisUrl)
	if err != nil {
		return nil, err
	}

	zLog, err := GetLogger(conf)
	if err != nil {
		return nil, fmt.Errorf("can't load logger: err=%v", err)
	}

	dateUC = uc.Date{}
	if err := dateUC.Init(); err != nil {
		return nil, fmt.Errorf("can't compaile regexp: err=%s", err.Error())
	}

	log = zLog
	redisClient = rClient
	botWrapper = h.NewBotWrapper(bot, zLog)

	handlers := map[h.HandlerKey]h.Handler{
		h.FileReportHandler:   reportHandler(),
		h.MappingHandler:      mappingHandler(),
		h.TransactionsHandler: transactionHandler(),
		h.TokenHandler:        tokenHandler(),
	}

	return h.NewChat(up, handlers, botWrapper), nil
}

func reportHandler() *h.FileReport {
	mappingRepo := repository.NewMapping(redisClient)
	tgRepo := repository.NewTelegram()
	fileReportUC := uc.NewFileReport(dateUC, mappingRepo, log, tgRepo)
	fileReportHandler := h.NewFileReport(fileReportUC, botWrapper)
	return fileReportHandler
}

func mappingHandler() *h.Mapping {
	mappingRepo := repository.NewMapping(redisClient)
	tgRepo := repository.NewTelegram()
	mappingUC := uc.NewMapping(mappingRepo, tgRepo)
	mappingHandler := h.NewMapping(mappingUC, botWrapper)
	return mappingHandler
}

func transactionHandler() *h.Transaction {
	tokenRepo := repository.NewToken(redisClient)
	apiRepo := repository.NewTransaction(log)
	mappingRepo := repository.NewMapping(redisClient)
	tokenUC := uc.NewToken(tokenRepo)
	apiUC := uc.NewTransaction(apiRepo, mappingRepo, log, dateUC)
	transactionHandler := h.NewTransaction(tokenUC, apiUC, botWrapper)
	return transactionHandler
}

func tokenHandler() *h.Token {
	tokenRepo := repository.NewToken(redisClient)
	tokenUC := uc.NewToken(tokenRepo)
	tokenHandler := h.NewToken(tokenUC, botWrapper)
	return tokenHandler
}
