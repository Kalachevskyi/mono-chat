// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/Kalachevskyi/mono-chat/app/adapters/mono"
	"github.com/Kalachevskyi/mono-chat/app/adapters/redis"
	telegram2 "github.com/Kalachevskyi/mono-chat/app/adapters/telegram"
	"github.com/Kalachevskyi/mono-chat/app/presetation/rest"
	"github.com/Kalachevskyi/mono-chat/app/presetation/telegram"
	"github.com/Kalachevskyi/mono-chat/app/usecases"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InjectReport(toolsWrapper ToolsWrapper) *telegram.FileReport {
	location := toolsWrapper.Loc
	date := usecases.NewDate(location)
	client := toolsWrapper.RedisClient
	mapping := redis.NewMapping(client)
	sugaredLogger := toolsWrapper.Log
	telegramTelegram := telegram2.NewTelegram()
	fileReport := usecases.NewFileReport(date, mapping, sugaredLogger, telegramTelegram)
	generic := redis.NewGeneric(client)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramFileReport := telegram.NewFileReport(fileReport, chatUser, botWrapper)
	return telegramFileReport
}

func InjectMapping(toolsWrapper ToolsWrapper) *telegram.Mapping {
	client := toolsWrapper.RedisClient
	mapping := redis.NewMapping(client)
	telegramTelegram := telegram2.NewTelegram()
	usecasesMapping := usecases.NewMapping(mapping, telegramTelegram)
	generic := redis.NewGeneric(client)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	sugaredLogger := toolsWrapper.Log
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramMapping := telegram.NewMapping(usecasesMapping, chatUser, botWrapper)
	return telegramMapping
}

func InjectTransaction(toolsWrapper ToolsWrapper) *telegram.Transaction {
	client := toolsWrapper.RedisClient
	generic := redis.NewGeneric(client)
	token := usecases.NewToken(generic)
	sugaredLogger := toolsWrapper.Log
	monoMono := mono.NewMono(sugaredLogger)
	mapping := redis.NewMapping(client)
	location := toolsWrapper.Loc
	date := usecases.NewDate(location)
	transaction := usecases.NewTransaction(monoMono, mapping, sugaredLogger, date)
	account := usecases.NewAccount(generic)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramTransaction := telegram.NewTransaction(token, transaction, account, chatUser, botWrapper)
	return telegramTransaction
}

func InjectToken(toolsWrapper ToolsWrapper) *telegram.Token {
	client := toolsWrapper.RedisClient
	generic := redis.NewGeneric(client)
	token := usecases.NewToken(generic)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	sugaredLogger := toolsWrapper.Log
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramToken := telegram.NewToken(token, chatUser, botWrapper)
	return telegramToken
}

func InjectClientInfo(toolsWrapper ToolsWrapper) *telegram.ClientInfo {
	client := toolsWrapper.RedisClient
	generic := redis.NewGeneric(client)
	token := usecases.NewToken(generic)
	sugaredLogger := toolsWrapper.Log
	monoMono := mono.NewMono(sugaredLogger)
	clientInfo := usecases.NewClientInfo(monoMono)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramClientInfo := telegram.NewClientInfo(token, clientInfo, chatUser, botWrapper)
	return telegramClientInfo
}

func InjectAccount(toolsWrapper ToolsWrapper) *telegram.Account {
	client := toolsWrapper.RedisClient
	generic := redis.NewGeneric(client)
	account := usecases.NewAccount(generic)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	sugaredLogger := toolsWrapper.Log
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramAccount := telegram.NewAccount(account, chatUser, botWrapper)
	return telegramAccount
}

func InjectUserChat(toolsWrapper ToolsWrapper) *telegram.ChatUser {
	client := toolsWrapper.RedisClient
	generic := redis.NewGeneric(client)
	chatUser := usecases.NewChatUser(generic)
	botAPI := toolsWrapper.Bot
	sugaredLogger := toolsWrapper.Log
	botWrapper := telegram.NewBotWrapper(botAPI, sugaredLogger)
	telegramChatUser := telegram.NewChatUser(chatUser, botWrapper)
	return telegramChatUser
}

func InjectTransactionRest(toolsWrapper ToolsWrapper) *rest.Transaction {
	sugaredLogger := toolsWrapper.Log
	monoMono := mono.NewMono(sugaredLogger)
	client := toolsWrapper.RedisClient
	mapping := redis.NewMapping(client)
	location := toolsWrapper.Loc
	date := usecases.NewDate(location)
	transaction := usecases.NewTransaction(monoMono, mapping, sugaredLogger, date)
	user := redis.NewUser(client)
	usecasesUser := usecases.NewUser(user)
	generic := redis.NewGeneric(client)
	account := usecases.NewAccount(generic)
	token := usecases.NewToken(generic)
	restTransaction := rest.NewTransaction(sugaredLogger, transaction, usecasesUser, account, token)
	return restTransaction
}

func InjectHTTPService(tw ToolsWrapper, port int) *rest.Service {
	sugaredLogger := tw.Log
	monoMono := mono.NewMono(sugaredLogger)
	client := tw.RedisClient
	mapping := redis.NewMapping(client)
	location := tw.Loc
	date := usecases.NewDate(location)
	transaction := usecases.NewTransaction(monoMono, mapping, sugaredLogger, date)
	user := redis.NewUser(client)
	usecasesUser := usecases.NewUser(user)
	generic := redis.NewGeneric(client)
	account := usecases.NewAccount(generic)
	token := usecases.NewToken(generic)
	restTransaction := rest.NewTransaction(sugaredLogger, transaction, usecasesUser, account, token)
	service := rest.NewService(restTransaction, port)
	return service
}

// wire.go:

var (
	fileReportUseCaseSet = wire.NewSet(usecases.NewFileReport, wire.Bind(new(telegram.CsvUC), new(*usecases.FileReport)))

	mappingUseCaseSet = wire.NewSet(usecases.NewMapping, wire.Bind(new(telegram.MappingUC), new(*usecases.Mapping)))

	transactionUseCaseSet = wire.NewSet(usecases.NewTransaction, wire.Bind(new(telegram.TransactionUC), new(*usecases.Transaction)), wire.Bind(new(rest.TransactionUC), new(*usecases.Transaction)))

	accountUseCaseSet = wire.NewSet(usecases.NewAccount, wire.Bind(new(telegram.AccountUC), new(*usecases.Account)), wire.Bind(new(rest.AccountUC), new(*usecases.Account)))

	chatUserUseCaseSet = wire.NewSet(usecases.NewChatUser, wire.Bind(new(telegram.ChatUserUC), new(*usecases.ChatUser)))

	userUseCaseSet = wire.NewSet(usecases.NewUser, wire.Bind(new(rest.UserUC), new(*usecases.User)))

	tokenUseCaseSet = wire.NewSet(usecases.NewToken, wire.Bind(new(telegram.TokenUC), new(*usecases.Token)), wire.Bind(new(rest.TokenUC), new(*usecases.Token)))

	clientInfoUseCaseSet = wire.NewSet(usecases.NewClientInfo, wire.Bind(new(telegram.ClientInfoUC), new(*usecases.ClientInfo)))

	mappingRepo = wire.NewSet(redis.NewMapping, wire.Bind(new(usecases.MappingRepo), new(*redis.Mapping)))

	genericRepo = wire.NewSet(redis.NewGeneric, wire.Bind(new(usecases.TokenRepo), new(*redis.Generic)), wire.Bind(new(usecases.AccountRepo), new(*redis.Generic)), wire.Bind(new(usecases.ChatUserRepo), new(*redis.Generic)))

	telegramRepo = wire.NewSet(telegram2.NewTelegram, wire.Bind(new(usecases.TelegramRepo), new(*telegram2.Telegram)))

	userRepo = wire.NewSet(redis.NewUser, wire.Bind(new(usecases.UserRepo), new(*redis.User)))

	monoRepo = wire.NewSet(mono.NewMono, wire.Bind(new(usecases.MonoRepo), new(*mono.Mono)))

	clientInfoRepo = wire.NewSet(mono.NewMono, wire.Bind(new(usecases.ClientInfoRepo), new(*mono.Mono)))

	apiLoggerBind      = wire.Bind(new(telegram.Logger), new(*zap.SugaredLogger))
	apiRestLoggerBind  = wire.Bind(new(rest.Logger), new(*zap.SugaredLogger))
	ucLoggerBind       = wire.Bind(new(usecases.Logger), new(*zap.SugaredLogger))
	telegramLoggerBind = wire.Bind(new(telegram2.Logger), new(*zap.SugaredLogger))
	monoLoggerBind     = wire.Bind(new(mono.Logger), new(*zap.SugaredLogger))

	toolsWrapperSet = wire.NewSet(wire.FieldsOf(new(ToolsWrapper), "Bot", "Log", "RedisClient", "Loc"))
)
