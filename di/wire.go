//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/Kalachevskyi/mono-chat/app/adapters/mono"
	ar "github.com/Kalachevskyi/mono-chat/app/adapters/redis"
	"github.com/Kalachevskyi/mono-chat/app/adapters/telegram"
	hr "github.com/Kalachevskyi/mono-chat/app/presetation/rest"
	h "github.com/Kalachevskyi/mono-chat/app/presetation/telegram"
	uc "github.com/Kalachevskyi/mono-chat/app/usecases"
)

var (
	fileReportUseCaseSet = wire.NewSet(
		uc.NewFileReport,
		wire.Bind(new(h.CsvUC), new(*uc.FileReport)),
	)

	mappingUseCaseSet = wire.NewSet(
		uc.NewMapping,
		wire.Bind(new(h.MappingUC), new(*uc.Mapping)),
	)

	transactionUseCaseSet = wire.NewSet(
		uc.NewTransaction,
		wire.Bind(new(h.TransactionUC), new(*uc.Transaction)),
		wire.Bind(new(hr.TransactionUC), new(*uc.Transaction)),
	)

	accountUseCaseSet = wire.NewSet(
		uc.NewAccount,
		wire.Bind(new(h.AccountUC), new(*uc.Account)),
		wire.Bind(new(hr.AccountUC), new(*uc.Account)),
	)

	chatUserUseCaseSet = wire.NewSet(
		uc.NewChatUser,
		wire.Bind(new(h.ChatUserUC), new(*uc.ChatUser)),
	)

	userUseCaseSet = wire.NewSet(
		uc.NewUser,
		wire.Bind(new(hr.UserUC), new(*uc.User)),
	)

	tokenUseCaseSet = wire.NewSet(
		uc.NewToken,
		wire.Bind(new(h.TokenUC), new(*uc.Token)),
		wire.Bind(new(hr.TokenUC), new(*uc.Token)),
	)

	clientInfoUseCaseSet = wire.NewSet(
		uc.NewClientInfo,
		wire.Bind(new(h.ClientInfoUC), new(*uc.ClientInfo)),
	)

	mappingRepo = wire.NewSet(
		ar.NewMapping,
		wire.Bind(new(uc.MappingRepo), new(*ar.Mapping)),
	)

	genericRepo = wire.NewSet(
		ar.NewGeneric,
		wire.Bind(new(uc.TokenRepo), new(*ar.Generic)),
		wire.Bind(new(uc.AccountRepo), new(*ar.Generic)),
		wire.Bind(new(uc.ChatUserRepo), new(*ar.Generic)),
	)

	telegramRepo = wire.NewSet(
		telegram.NewTelegram,
		wire.Bind(new(uc.TelegramRepo), new(*telegram.Telegram)),
	)

	userRepo = wire.NewSet(
		ar.NewUser,
		wire.Bind(new(uc.UserRepo), new(*ar.User)),
	)

	monoRepo = wire.NewSet(
		mono.NewMono,
		wire.Bind(new(uc.MonoRepo), new(*mono.Mono)),
	)

	clientInfoRepo = wire.NewSet(
		mono.NewMono,
		wire.Bind(new(uc.ClientInfoRepo), new(*mono.Mono)),
	)

	apiLoggerBind      = wire.Bind(new(h.Logger), new(*zap.SugaredLogger))
	apiRestLoggerBind  = wire.Bind(new(hr.Logger), new(*zap.SugaredLogger))
	ucLoggerBind       = wire.Bind(new(uc.Logger), new(*zap.SugaredLogger))
	telegramLoggerBind = wire.Bind(new(telegram.Logger), new(*zap.SugaredLogger))
	monoLoggerBind     = wire.Bind(new(mono.Logger), new(*zap.SugaredLogger))

	toolsWrapperSet = wire.NewSet(
		wire.FieldsOf(new(ToolsWrapper), "Bot", "Log", "RedisClient", "Loc"),
	)
)

func InjectReport(ToolsWrapper) *h.FileReport {
	wire.Build(
		toolsWrapperSet,
		h.NewFileReport,
		uc.NewDate,
		chatUserUseCaseSet,
		genericRepo,
		mappingRepo,
		telegramRepo,
		h.NewBotWrapper,
		fileReportUseCaseSet,
		ucLoggerBind,
		apiLoggerBind,
	)
	return nil
}

func InjectMapping(ToolsWrapper) *h.Mapping {
	wire.Build(
		toolsWrapperSet,
		h.NewMapping,
		h.NewBotWrapper,
		mappingUseCaseSet,
		chatUserUseCaseSet,
		genericRepo,
		mappingRepo,
		telegramRepo,
		apiLoggerBind,
	)
	return nil
}

func InjectTransaction(ToolsWrapper) *h.Transaction {
	wire.Build(
		h.NewTransaction,
		toolsWrapperSet,
		tokenUseCaseSet,
		chatUserUseCaseSet,
		genericRepo,
		transactionUseCaseSet,
		accountUseCaseSet,
		mappingRepo,
		uc.NewDate,
		monoRepo,
		h.NewBotWrapper,
		apiLoggerBind,
		monoLoggerBind,
		ucLoggerBind,
	)
	return nil
}

func InjectToken(ToolsWrapper) *h.Token {
	wire.Build(
		h.NewToken,
		toolsWrapperSet,
		tokenUseCaseSet,
		chatUserUseCaseSet,
		genericRepo,
		h.NewBotWrapper,
		apiLoggerBind,
	)
	return nil
}

func InjectClientInfo(ToolsWrapper) *h.ClientInfo {
	wire.Build(
		h.NewClientInfo,
		toolsWrapperSet,
		tokenUseCaseSet,
		clientInfoUseCaseSet,
		chatUserUseCaseSet,
		genericRepo,
		clientInfoRepo,
		monoLoggerBind,
		h.NewBotWrapper,
		apiLoggerBind,
	)
	return nil
}

func InjectAccount(ToolsWrapper) *h.Account {
	wire.Build(
		h.NewAccount,
		toolsWrapperSet,
		accountUseCaseSet,
		chatUserUseCaseSet,
		genericRepo,
		h.NewBotWrapper,
		apiLoggerBind,
	)
	return nil
}

func InjectUserChat(ToolsWrapper) *h.ChatUser {
	wire.Build(
		h.NewChatUser,
		toolsWrapperSet,
		chatUserUseCaseSet,
		genericRepo,
		h.NewBotWrapper,
		apiLoggerBind,
	)
	return nil
}

func InjectTransactionRest(ToolsWrapper) *hr.Transaction {
	wire.Build(
		hr.NewTransaction,
		toolsWrapperSet,
		tokenUseCaseSet,
		genericRepo,
		transactionUseCaseSet,
		accountUseCaseSet,
		userUseCaseSet,
		mappingRepo,
		uc.NewDate,
		monoRepo,
		userRepo,
		ucLoggerBind,
		apiRestLoggerBind,
		monoLoggerBind,
	)
	return nil
}

func InjectHTTPService(tw ToolsWrapper, port int) *hr.Service {
	wire.Build(
		hr.NewService,
		hr.NewTransaction,
		toolsWrapperSet,
		tokenUseCaseSet,
		genericRepo,
		transactionUseCaseSet,
		accountUseCaseSet,
		userUseCaseSet,
		mappingRepo,
		uc.NewDate,
		monoRepo,
		userRepo,
		ucLoggerBind,
		apiRestLoggerBind,
		monoLoggerBind,
	)
	return nil
}
