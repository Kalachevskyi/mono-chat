package di

import (
	"time"

	"github.com/go-redis/redis"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

// ToolsWrapper represents tools wrapper.
type ToolsWrapper struct {
	Log         *zap.SugaredLogger
	RedisClient *redis.Client
	Loc         *time.Location
	Bot         *tg.BotAPI
}
