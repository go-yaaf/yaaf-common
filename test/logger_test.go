// Logger tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	skipCI(t)
	logger.EnableJsonFormat(false)
	logger.Init()

	logger.Debug("This should have an ISO8601 based time stamp debug")
	logger.Info("This should have an ISO8601 based time stamp info")
	logger.Warn("This should have an ISO8601 based time stamp warning")
	logger.Error("This should have an ISO8601 based time stamp error")
}

func TestProductionLogger(t *testing.T) {
	skipCI(t)
	logger.Debug("debug message")
	time.Sleep(time.Second)
	logger.Info("info message")
	time.Sleep(time.Second)
	logger.Warn("warning message")
	time.Sleep(time.Second)
	logger.Error("error message")
}

func TestZapLogger(t *testing.T) {

	skipCI(t)

	cfg := zap.Config{
		Encoding:    "console",
		OutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			TimeKey:     "time",
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
		},
	}

	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	fmt.Printf("*** Using a standard encoder\n")

	logger2, _ := cfg.Build()
	logger2.Debug("This should have an ISO8601 based time stamp debug")
	logger2.Info("This should have an ISO8601 based time stamp info")
	logger2.Warn("This should have an ISO8601 based time stamp warning")
	logger2.Error("This should have an ISO8601 based time stamp error")

	fmt.Printf("*** Using a custom encoder\n")

	cfg.EncoderConfig.EncodeLevel = CustomLevelEncoder
	logger2, _ = cfg.Build()
	logger2.Info("This should have a interesting level name")
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01/02  15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
