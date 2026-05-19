package logger

import (
	"os"
	"sync"

	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapSinLogger *zap.SugaredLogger
var once sync.Once

type ZapLogger struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewZapLogger(cfg *config.Config) *ZapLogger {
	logger := &ZapLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (l *ZapLogger) Init() {
	once.Do(func() {
		zapConfig := zap.NewProductionEncoderConfig()
		zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleConfig := zap.NewDevelopmentEncoderConfig()
		consoleConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		var core zapcore.Core
		if l.cfg.Server.RunMode != "release" {
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(consoleConfig),
				zapcore.AddSync(os.Stdout),
				zap.DebugLevel,
			)
		} else {

		}

		zapSinLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1)).Sugar()
	})

	l.logger = zapSinLogger
}

func (l *ZapLogger) GetLogger() *zap.Logger {
	return l.logger.Desugar()
}

func (l *ZapLogger) Info(msg string, args ...any) {
	l.logger.Infow(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...any) {
	l.logger.Errorw(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...any) {
	l.logger.Warnw(msg, args...)
}

func (l *ZapLogger) Debug(msg string, args ...any) {
	l.logger.Debugw(msg, args...)
}

func (l *ZapLogger) Fatal(msg string, args ...any) {
	l.logger.Fatalw(msg, args...)
}

func (l *ZapLogger) Sync() {
	err := l.logger.Sync()
	if err != nil {
		return
	}
}
