/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

// Package logger is logger package for smart contract sdk.
package logger

import (
	"fmt"

	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

const (
	empty         = ""
	logTypePretty = "pretty"
	logTypeJSON   = "json"
	moduleNameKey = "module"
	loggerNameKey = "name"

	defaultLevel      = zerolog.InfoLevel
	defaultLogType    = "json"
	defaultTimeformat = "2006-01-02 15:04:05.000 -0700"
	initialized       = iota
	notInitialized
)

const callerWithSkipFrameCount = 3

// LogCfg log config.
type LogCfg struct {
	Type         string
	Level        string
	ModuleLevels map[string]string
}

var wienerLogger *WienerLogger
var moduleLevels map[string]zerolog.Level

type observer func(writer io.Writer, level zerolog.Level)

var hasInit int32 = notInitialized
var observers []observer
var obLock sync.Mutex

/// nolint
func init() {
	writer, err := getLogWriter(defaultLogType)
	if err != nil {
		return
	}
	logger := zerolog.New(writer).
		Level(defaultLevel).
		With().
		Timestamp().
		CallerWithSkipFrameCount(callerWithSkipFrameCount).
		Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	wienerLogger = &WienerLogger{logger: &logger}
}

// Init init for log.
func Init(config *LogCfg) error {
	writer, err := getLogWriter(config.Type)
	if err != nil {
		return err
	}
	level, err := zerolog.ParseLevel(strings.ToLower(config.Level))
	if err != nil {
		return fmt.Errorf("unknown log level: %s", config.Level)
	}
	logger := wienerLogger.logger.
		Output(writer).
		Level(level)
	wienerLogger.logger = &logger

	moduleLevels = make(map[string]zerolog.Level)
	for module, levelStr := range config.ModuleLevels {
		moduleLevel, err := zerolog.ParseLevel(strings.ToLower(levelStr))
		if err != nil {
			return fmt.Errorf("unknown module log level: module - %s, level - %s", module, levelStr)
		}
		moduleLevels[module] = moduleLevel
	}
	atomic.StoreInt32(&hasInit, initialized)
	for _, ob := range observers {
		ob(writer, level)
	}
	observers = nil
	return nil
}

// GetDefaultLogger get default log handler.
func GetDefaultLogger() *WienerLogger {
	return wienerLogger
}

// GetLogger get log handler by log name.
func GetLogger(loggerName string) *WienerLogger {
	return GetModuleLogger(empty, loggerName)
}

// GetModuleLogger different modules might have different log level, which can be set
// in config file.
func GetModuleLogger(moduleName, loggerName string) *WienerLogger {
	subLogger := wienerLogger.logger.
		With().
		Str(moduleNameKey, moduleName).
		Str(loggerNameKey, loggerName).
		Logger()
	wienerLogger := &WienerLogger{logger: &subLogger}

	init := atomic.LoadInt32(&hasInit)
	if init == notInitialized {
		addObserver(func(writer io.Writer, level zerolog.Level) {
			newLevel := level
			moduleLevel, ok := moduleLevels[moduleName]
			if ok {
				newLevel = moduleLevel
			}
			logger := subLogger.Output(writer).Level(newLevel)
			wienerLogger.logger = &logger
		})
	} else {
		moduleLevel, ok := moduleLevels[moduleName]
		if ok {
			logger := subLogger.Level(moduleLevel)
			wienerLogger.logger = &logger
		}
	}
	return wienerLogger
}

func getLogWriter(logType string) (io.Writer, error) {
	switch strings.ToLower(logType) {
	case logTypePretty:
		return &Writer{out: os.Stderr, timeFormat: defaultTimeformat}, nil
	case logTypeJSON:
		return os.Stderr, nil
	case empty:
		return &Writer{out: os.Stderr, timeFormat: defaultTimeformat}, nil
	}
	return nil, fmt.Errorf("unknown log type: %s", logType)
}

func addObserver(ob observer) {
	obLock.Lock()
	observers = append(observers, ob)
	obLock.Unlock()
}
