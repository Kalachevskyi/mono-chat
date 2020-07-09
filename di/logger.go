package di

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger //nolint:gochecknoglobals
	loggerOnce sync.Once          //nolint:gochecknoglobals
)

// Logger - initialize "zap" logger, returns sugared instance of logger
func Logger(debug bool, encodingLog string) (*zap.SugaredLogger, error) {
	var err error

	loggerOnce.Do(func() {
		zConf := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: debug,
			Sampling: &zap.SamplingConfig{
				Initial:    100, //nolint:gomnd
				Thereafter: 100, //nolint:gomnd
			},
			Encoding:          encodingLog,
			EncoderConfig:     zap.NewProductionEncoderConfig(),
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
			DisableCaller:     true,
			DisableStacktrace: true,
		}
		var zLog *zap.Logger
		zLog, err = zConf.Build()
		if err != nil {
			return
		}
		logger = zLog.Sugar()
	})

	return logger, err
}
