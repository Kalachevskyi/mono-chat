package infrastructure

import (
	"github.com/Kalachevskyi/mono-chat/config"
	"go.uber.org/zap"
)

func GetLogger(c config.Config) (*zap.SugaredLogger, error) {
	zConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: c.Debug,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          c.EncodingLog,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}
	zLog, err := zConf.Build()
	if err != nil {
		return nil, err
	}

	sugar := zLog.Sugar()
	return sugar, nil
}
