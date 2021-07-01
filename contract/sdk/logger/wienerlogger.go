/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

package logger

import (
	"github.com/rs/zerolog"
)

// WienerLogger wienerchain log.
type WienerLogger struct {
	logger *zerolog.Logger
}

// Trace the trace level log.
func (l *WienerLogger) Trace(msg string) {
	l.logger.Trace().Msg(msg)
}

// Tracef the tracef level log.
func (l *WienerLogger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

// Debug the debug level log.
func (l *WienerLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf the debug level log.
func (l *WienerLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info the info level log.
func (l *WienerLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof the infof level log.
func (l *WienerLogger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn the warn level log.
func (l *WienerLogger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf the warnf level log.
func (l *WienerLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error the error level log.
func (l *WienerLogger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf the errorf level log.
func (l *WienerLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatal the fatal level log.
func (l *WienerLogger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf the fatalf level log.
func (l *WienerLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// Panic the panic level log.
func (l *WienerLogger) Panic(msg string) {
	l.logger.Panic().Msg(msg)
}

// Panicf the panicf level log.
func (l *WienerLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}
