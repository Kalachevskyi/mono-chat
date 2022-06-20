package di

import (
	"go.uber.org/zap"
)

// Logger - initialize "zap" logger, returns sugared instance of logger.
func Logger(debug bool, encodingLog string) (*zap.SugaredLogger, error) {
	zConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: debug,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          encodingLog,
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

	return zLog.Sugar(), err
}
