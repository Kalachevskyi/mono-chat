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

// Package framework is an application layer for initializing app components
package framework

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger //nolint:gochecknoglobals
	loggerOnce sync.Once          //nolint:gochecknoglobals
)

// GetLogger - initialize "zap" logger, returns sugared instance of logger
func GetLogger(debug bool, encodingLog string) (*zap.SugaredLogger, error) {
	var err error
	loggerOnce.Do(func() {
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
		var zLog *zap.Logger
		zLog, err = zConf.Build()
		if err != nil {
			return
		}
		logger = zLog.Sugar()
	})

	return logger, err
}
