package service

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type TimerManager struct {
	logger *zap.Logger
}

func newTimerManager(logPath string) (*TimerManager, error) {
	m := TimerManager{}
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{logPath},
		ErrorOutputPaths: []string{logPath},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	m.logger = log
	return &m, nil
}

func (m *TimerManager) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		defer func() {
			m.logger.Info("request was executed", zap.String("URI", r.RequestURI),
				zap.String("took", time.Since(st).String()))
		}()
		next.ServeHTTP(w, r)
	})
}
