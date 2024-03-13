package service

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		config *Config
		level  log.Level
		err    bool
	}{
		{config: &Config{LogLevel: "panic"}, level: log.PanicLevel},
		{config: &Config{LogLevel: "fatal"}, level: log.FatalLevel},
		{config: &Config{LogLevel: "error"}, level: log.ErrorLevel},
		{config: &Config{LogLevel: "warn"}, level: log.WarnLevel},
		{config: &Config{LogLevel: "warning"}, level: log.WarnLevel},
		{config: &Config{LogLevel: "info"}, level: log.InfoLevel},
		{config: &Config{LogLevel: "debug"}, level: log.DebugLevel},
		{config: &Config{LogLevel: "trace"}, level: log.TraceLevel},
		{config: &Config{LogLevel: "invalid"}, err: true},
	}

	for _, tt := range tests {
		if tt.err {
			assert.Panics(t, func() { tt.config.mustParseLogLevel() })
		} else {
			assert.Equal(t, tt.level, tt.config.mustParseLogLevel())
		}
	}
}
